package kafwifi

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"os"
)

func formatBytesLength(length int64) string {
	if length < 1024*1024 {
		return fmt.Sprintf("%.2f K", float32(length)/(1024))
	} else {
		return fmt.Sprintf("%.2f M", float32(length)/(1024*1024))
	}
}

func isExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func getClientID() string {
	clientID := uuid.NewV4().String()
	config, err := os.UserConfigDir()
	if err != nil {
		return clientID
	}
	filepath := fmt.Sprintf("%s/kaf-wifi/config", config)
	if exist, _ := isExists(filepath); exist {
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
