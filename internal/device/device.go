package device

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"
)

type Device struct {
	IP   string
	Port string
}

func (d *Device) SendCommand(command Command) error {
	address := net.JoinHostPort(d.IP, d.Port)
	conn, err := net.DialTimeout("tcp", address, 3*time.Second)
	if err != nil {
		return fmt.Errorf("не удалось подключиться к %s: %w", address, err)
	}

	request := fmt.Sprintf(`{"command":"%s"}`, command.GetCommand())
	_, err = fmt.Fprintf(conn, request)
	if err != nil {
		return fmt.Errorf("ошибка при отправке команды: %w", err)
	}

	response, err := io.ReadAll(conn)
	if err != nil {
		return fmt.Errorf("ошибка при получении ответа: %w", err)
	}

	err = json.Unmarshal(response[:len(response)-1], command.GetResponse())
	if err != nil {
		return fmt.Errorf("не удалось разобрать JSON-ответ: %w", err)
	}

	err = conn.Close()
	if err != nil {
		return fmt.Errorf("error closing connection: %w", err)
	}

	return nil
}
