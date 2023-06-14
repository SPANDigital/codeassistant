package vertexai

import (
	"os/exec"
	"strings"
)

func generateAccessToken() (string, error) {
	out, err := exec.Command("/Users/richardwooding/Downloads/google-cloud-sdk/bin/gcloud", "auth", "print-access-token").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
