package systemd

import (
	"bytes"
	"os/exec"
)

type serviceError string

func (e serviceError) Error() string { return string(e) }

// RestartSystemdService restarts the service and return status code
// String -> int, err
// RestartSystemdService(nginx) -> 0 		service restarted, no error
// RestartSystemdService(unknown) -> 1		response code 1, error - no such service
func RestartSystemdService(service string) error {

	if service == "" {
		return serviceError("service name can't be empty")
	}

	cmd := exec.Command("systemctl", "restart", service+".service")

	if bytes.HasSuffix([]byte(service), []byte(".service")) {
		cmd = exec.Command("systemctl", "restart", service)
	}

	return cmd.Run()
}
