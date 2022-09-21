package main

import (
	"os"

	"github.com/csaldiasdev/distws/internal/util"
)

func main() {
	ip, err := util.GetLocalIp()
	if err != nil {
		os.Exit(1)
	}
	os.Stdout.Write([]byte(ip))
}
