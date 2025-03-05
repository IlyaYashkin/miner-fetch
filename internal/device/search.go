package device

import (
	"bufio"
	"context"
	"os/exec"
	"regexp"
	"strings"
)

func GetRustScanScanner() *RustScanScanner {
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
