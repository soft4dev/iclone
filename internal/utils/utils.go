package utils

import (
	"fmt"
	"os/exec"
	"strings"
)

func CheckGitInstalled() error {
	cmd := exec.Command("git", "--version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error: git is not installed")
	}
	return nil
}

func ContainsStringInStringSlice(slice []string, s string) bool {
	for _, v := range slice {
		if strings.EqualFold(v, s) {
			return true
		}
	}
	return false
}
