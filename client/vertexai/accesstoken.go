package vertexai

import (
	"os/exec"
	"strings"
)

func generateAccessToken(gcloudBinaryPath string) (string, error) {
	out, err := exec.Command(gcloudBinaryPath, "auth", "print-access-token").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
