package device

import (
	"bufio"
	"context"
	"fmt"
	"github.com/Ullaakut/nmap/v3"
	"log"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

type Scanner interface {
	Scan(ctx context.Context) ([]Device, error)
}

func GetRustScanScanner() Scanner {
	return &RustScanScanner{}
}

type RustScanScanner struct{}

func (r RustScanScanner) Scan(ctx context.Context) ([]Device, error) {
	var devices []Device

	addresses := "192.168.0.0/24"
	ports := "4028"

	cmd := exec.CommandContext(ctx, "rustscan", "-g", "-a", addresses, "-p", ports)

	output, err := cmd.Output()
	if err != nil {
		return devices, err
	}

	re := regexp.MustCompile(`(\d+\.\d+\.\d+\.\d+) -> \[(.+)]`)

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindStringSubmatch(line)
		if len(matches) == 3 {
			ip := matches[1]
			ports := strings.Split(matches[2], ",")

			devices = append(devices, Device{IP: ip, Port: ports[0]})
		}
	}

	if err := scanner.Err(); err != nil {
		return devices, err
	}

	return devices, nil
}

type NmapScanner struct{}

func GetNmapScanner() Scanner {
	return &NmapScanner{}
}

func (n *NmapScanner) Scan(ctx context.Context) ([]Device, error) {
	ctxt, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	scanner, err := nmap.NewScanner(
		ctxt,
		nmap.WithTargets("192.168.0.0/24"),
		nmap.WithPorts("4028"),
	)
	if err != nil {
		return []Device{}, err
	}

	progress := make(chan float32)
	doneAsync := make(chan error)

	scanner.Async(doneAsync)
	scanner.Progress(progress)

	result, warnings, err := scanner.Run()

L:
	for {
		select {
		case value := <-progress:
			log.Printf("Progress %.2f%%", value)
		case err := <-doneAsync:
			if err != nil {
				log.Println("doneAsync:", err)
			}
			break L
		}
	}

	if len(*warnings) > 0 {
		log.Printf("run finished with warnings: %s\n", *warnings)
	}
	if err != nil {
		return []Device{}, err
	}

	var devices []Device

	for _, host := range result.Hosts {
		if len(host.Ports) == 0 || len(host.Addresses) == 0 {
			continue
		}

		for _, port := range host.Ports {
			if port.State.State != "open" {
				continue
			}

			devices = append(devices, Device{IP: host.Addresses[0].Addr, Port: fmt.Sprintf("%d", port.ID)})
		}
	}

	return devices, nil
}

type MockScanner struct{}

func GetMockScanner() Scanner {
	return &MockScanner{}
}

func (m *MockScanner) Scan(ctx context.Context) ([]Device, error) {
	return []Device{
		{IP: "192.168.0.19", Port: "4028"},
		{IP: "192.168.0.20", Port: "4028"},
	}, nil
}
