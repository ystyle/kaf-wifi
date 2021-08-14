package kafwifi

import (
	"fmt"
	"github.com/ystyle/google-analytics"
	"math/rand"
	"runtime"
	"time"
)

func Analytics(secret, measurement string) {
	if secret == "" || measurement == "" {
		return
	}
	t := time.Now().Unix()
	analytics.SetKeys(secret, measurement) // // required
	payload := analytics.Payload{
		ClientID: fmt.Sprintf("%d.%d", rand.Int31(), t), // required
		UserID:   getClientID(),
		Events: []analytics.Event{
			{
				Name: "kaf_wifi", // required
				Params: map[string]interface{}{
					"os":      runtime.GOOS,
					"arch":    runtime.GOARCH,
					"version": version,
				},
			},
		},
	}
	analytics.Send(payload)
}
