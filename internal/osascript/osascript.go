package osascript

import (
	"os/exec"
	"strings"
)

// Run executes an AppleScript and returns the output
func Run(script string) (string, error) {
	cmd := exec.Command("osascript", "-e", script)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// RunJS executes JavaScript for Automation (JXA) and returns the output
func RunJS(script string) (string, error) {
	cmd := exec.Command("osascript", "-l", "JavaScript", "-e", script)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
