package main

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
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

func IsExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func GetClientID() string {
	clientID := uuid.NewV4().String()
	config, err := os.UserConfigDir()
	if err != nil {
		return clientID
	}
	filepath := fmt.Sprintf("%s/kaf-wifi/config", config)
	if exist, _ := IsExists(filepath); exist {
		bs, err := ioutil.ReadFile(filepath)
		if err != nil {
			return clientID
		}
		clientID = string(bs)
	} else {
		err := os.MkdirAll(fmt.Sprintf("%s/kaf-wifi", config), 0700)
		if err != nil {
			return clientID
		}
		_ = os.WriteFile(filepath, []byte(clientID), 0700)
	}
	return clientID
}
