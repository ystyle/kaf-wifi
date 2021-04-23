package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"time"
)

var (
	secret      string
	measurement string
)

func analytics() {
	query := url.Values{}
	query.Add("api_secret", secret)
	query.Add("measurement_id", measurement)
	uri := fmt.Sprintf("https://www.google-analytics.com/mp/collect?%s", query.Encode())
	t := time.Now().Unix()
	params := map[string]interface{}{
		"client_id": fmt.Sprintf("%d.%d", rand.Int31(), t),
		"user_id":   GetClientID(),
		"events": []map[string]interface{}{
			{
				"name": "kaf_wifi",
				"params": map[string]string{
					"os":      runtime.GOOS,
					"arch":    runtime.GOARCH,
					"version": version,
				},
			},
		},
	}
	bs, _ := json.Marshal(params)
	payload := bytes.NewReader(bs)
	_, _ = http.Post(uri, "application/json", payload)
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
