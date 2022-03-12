package clihelpers

import (
	"bytes"
	"log"
	"os/exec"
	"strings"
)

func RunCliCommand(name string, params ...string) (string, string, error) {
	cmd := exec.Command(name, params...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func GetYoutubeDownloadLink(youtubeLink string) (string, bool) {
	var link string
	var err error
	// retries
	for i := 0; i < 3; i++ {
		link, _, err = RunCliCommand("youtube-dl", "-x", "-g", youtubeLink)
		if err == nil {
			break
		}
	}
	exists := true
	if err != nil {
		exists = false
		log.Println("error querying youtube-dl:", err)
	}

	return strings.TrimSpace(link), exists
}

type FfmpegMetadata struct {
	Title  string
	Artist string
	Album  string
	Image  string
}

func FfmpegConvert(filepathIn, filepathOut string, metadata FfmpegMetadata) error {
	args := make([]string, 0, 10)
	args = append(args, "-y", "-i", filepathIn)

	if len(metadata.Image) > 0 {
		args = append(args, "-i", metadata.Image)
		args = append(args, "-map", "0:0", "-map", "1:0")
		args = append(args, "-metadata:s:v", "title=Album cover")
		args = append(args, "-metadata:s:v", "comment=Cover (Front)")
	}

	args = append(args, "-id3v2_version", "3")

	args = append(args, "-metadata:s:a", "title="+metadata.Title)
	args = append(args, "-metadata:s:a", "artist="+metadata.Artist)
	args = append(args, "-metadata:s:a", "album="+metadata.Album)

	args = append(args, filepathOut)

	var err error
	// convertation retries
	for i := 0; i < 3; i++ {
		_, _, err = RunCliCommand("ffmpeg", args...)
		if err == nil {
			break
		}
	}

	return err
}
