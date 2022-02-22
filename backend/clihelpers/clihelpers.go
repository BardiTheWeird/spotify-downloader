package clihelpers

import (
	"bytes"
	"log"
	"os/exec"
	"strings"
)

func RunCliCommand(name string, params ...string) (string, string, error) {
	// get the download link with the power of youtube-dl
	cmd := exec.Command(name, params...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func GetYoutubeDownloadLink(youtubeLink string) (string, bool) {
	link, _, err := RunCliCommand("youtube-dl", "-x", "-g", youtubeLink)
	exists := true
	if err != nil {
		exists = false
		log.Println("error querying youtube-dl:", err)
	}

	return strings.TrimSpace(link), exists
}
