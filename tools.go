package main

import (
	"fmt"
	"os"
	"os/exec"
)

func FormatBytesLength(length int64) string {
	if length < 1024*1024 {
		return fmt.Sprintf("%d K", length/(1024))
	} else {
		return fmt.Sprintf("%d M", length/(1024*1024))
	}
}

func Run(dir, command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
