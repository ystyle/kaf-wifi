GOOS=linux GOARCH=arm go build -ldflags "-X main.secret=$secret -X main.measurement=$measurement" kaf-wifi
