package main

import (
	"log"

	"github.com/sullivtr/k8s_platform/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalf("[ERROR] %v", err)
	}
}
