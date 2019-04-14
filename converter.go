package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

// StorageDirPath is where to save a downloaded file
var StorageDirPath string

// ConvertFile downloads and converts file
func ConvertFile(url string) (string, error) {

	// firstly, download the video file
	tmpFileName, err := fetchRemote(url)
	if err != nil {
		return tmpFileName, err
	}

	// then, convert it
	if err := convert(tmpFileName); err != nil {
		return tmpFileName, err
	}

	return tmpFileName, nil
}

func fetchRemote(fileURL string) (string, error) {

	fileName := uuid.Must(uuid.NewV4()).String() // generate new name avoiding collision

	if strings.HasSuffix(fileURL, ".gifv") {
		fileURL = strings.Replace(fileURL, ".gifv", ".mp4", 1)
	}

	fileToConvert := fmt.Sprintf("%s/%s.mp4", StorageDirPath, fileName)
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

	return fileName, nil
}

func convert(fileName string) error {

	// Convert movie to gif
	fileToConvert := fmt.Sprintf("%s/%s.mp4", StorageDirPath, fileName)
	outputFile := fmt.Sprintf("%s/%s.gif", StorageDirPath, fileName)
	ffmpeg := exec.Command("ffmpeg", "-i", fileToConvert, "-pix_fmt", "rgb24", "-vf", "scale=300:-1", "-f", "gif", outputFile)

	var ffmpegErr bytes.Buffer
	ffmpeg.Stderr = &ffmpegErr

	err := ffmpeg.Run()
	if err != nil {
		log.Println("Can't convert to GIF: " + err.Error())
		return errors.New(fmt.Sprint(err) + ": " + ffmpegErr.String())
	}

	// Optimize gif
	sickle := exec.Command("gifsicle", "--careful", "-O3", "--batch", outputFile)

	var sicklekErr bytes.Buffer
	sickle.Stderr = &sicklekErr

	err = sickle.Run()
	if err != nil {
		log.Println("Can't optimyze GIF: " + err.Error())
		return errors.New(fmt.Sprint(err) + ": " + sicklekErr.String())
	}

	return nil
}

// CleanUp remove all the generated and downloaded files
func CleanUp(fileName string) error {
	fileToConvert := fmt.Sprintf("%s/%s.mp4", StorageDirPath, fileName)
	outputFile := fmt.Sprintf("%s/%s.gif", StorageDirPath, fileName)

	if err := os.Remove(fileToConvert); err != nil {
		log.Printf("Can't delete %s file, error: %s", fileToConvert, err.Error())
		return err
	}

	if err := os.Remove(outputFile); err != nil {
		log.Printf("Can't delete %s file, error: %s", outputFile, err.Error())
		return err
	}

	return nil
}
