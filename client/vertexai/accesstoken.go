package vertexai

import (
	"github.com/spf13/viper"
	"os/exec"
	"strings"
)

func generateAccessToken() (string, error) {
	out, err := exec.Command(viper.GetString("gcloudBinary"), "auth", "print-access-token").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
