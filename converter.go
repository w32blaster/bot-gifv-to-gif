package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

func ConvertFile(url string) error {

	tmpFilePath, err := fetchRemote(url)
	if err != nil {
		return err
	}
	convert(tmpFilePath)

	return nil
}

func fetchRemote(fileURL string) (string, error) {

	url, err := url.Parse(fileURL)
	if err != nil {
		return "", err
	}

	fileExt := path.Ext(url.Path)

	// Gifv is a container for mp4
	if fileExt == ".gifv" {
		fileURL = strings.Replace(fileURL, ".gifv", ".mp4", -1)
	}

	fileToConvert := "/tmp/temp_file_to_convert" + fileExt // REPLACE WITH TIMESTAMP
	temp, err := os.Create(fileToConvert)
	defer temp.Close()

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", fileURL, nil)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	_, err = io.Copy(temp, resp.Body)
	if err != nil {
		return "", err
	}

	return fileToConvert, nil
}

func convert(tmpFilePath string) error {

	// Convert movie to gif
	outputImage := "/tmp/temp_file_to_convert2.gif"
	ffmpeg := exec.Command("ffmpeg", "-i", tmpFilePath, "-pix_fmt", "rgb24", "-vf", "scale=300:-1", "-f", "gif", outputImage)

	var ffmpegErr bytes.Buffer
	ffmpeg.Stderr = &ffmpegErr

	err := ffmpeg.Run()
	if err != nil {
		return errors.New(fmt.Sprint(err) + ": " + ffmpegErr.String())
	}

	// Optimize gif
	sickle := exec.Command("gifsicle", "--careful", "-O3", "--batch", outputImage)

	var sicklekErr bytes.Buffer
	sickle.Stderr = &sicklekErr

	err = sickle.Run()
	if err != nil {
		return errors.New(fmt.Sprint(err) + ": " + sicklekErr.String())
	}

	return nil
}
