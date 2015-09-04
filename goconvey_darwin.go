package main

import (
	"fmt"
	"log"
	"os/exec"
)

func launchBrowser(host string, port int) {
	log.Printf("Launching browser on %s:%d", host, port)

	url := fmt.Sprintf("http://%s:%d", host, port)
	cmd := exec.Command("open", url)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	log.Println(string(output))
}
