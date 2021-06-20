package main

import (
	kafwifi "kaf-wifi"
	"os"
)

var (
	analyze     = true
	secret      string
	measurement string
)

func main() {
	if len(os.Args) == 3 && os.Args[2] == "-noanalyze" {
		analyze = false
	}
	if analyze {
		go kafwifi.Analytics(secret, measurement)
	}
	kafwifi.Start()
}
