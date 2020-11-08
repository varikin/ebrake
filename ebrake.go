package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

type Encoder struct {
	config *Config
	Source string
	Target string
}

type Video struct {
	source string
	target string
}

// EncodeFiles encodes all the video files in the source directory and places them
// in the destination directory.
func (encoder *Encoder) EncodeFiles() error {

	info, err := os.Stat(encoder.Target)
	if err != nil {
		fmt.Println("Target directory does not exist; attempting to create it.")
		if err2 := os.MkdirAll(encoder.Target, 0666); err2 != nil {
			return errors.Wrap(err, "Target directory does not exist")
		}
	} else if !info.IsDir() {
		return errors.New(encoder.Target + " is not a directory")
	}

	videoFiles, err := encoder.getVideoFiles()

	if err != nil {
		return err
	}

	var videos []Video
	for _, videoFile := range videoFiles {
		targetFile, err := encoder.getDestinationPath(videoFile)
		if err != nil {
			return err
		}

		// Filter out targetFiles that exist
		// TODO make this an option, like --force
		exists, err := fileExists(targetFile)
		if err != nil {
			return err
		}
		if exists {
			fmt.Println("Target file already exists, skipping: " + targetFile)
		} else {
			videos = append(videos, Video{
				source: videoFile,
				target: targetFile,
			})
		}
	}
	if len(videos) == 0 {
		fmt.Println("Did not find any videos re-encode.")
		return nil
	}

	// Encode each file
	options := strings.Fields(encoder.config.HandBrakeOptions)
	for _, video := range videos {
		args := append(options, "-i", video.source, "-o", video.target)
		cmd := exec.Command(encoder.config.HandBrakeCommand, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		if err != nil {
			return errors.Wrap(err, "failed to encode video: "+video.source)
		}
	}

	return nil
}

func fileExists(name string) (bool, error) {
	_, err := os.Stat(name)

	// No error means the file exists
	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, errors.Wrap(err, "Unable to determine state of file: "+name)
	}
}

// getVideoFiles finds and returns the list of video files in the given directory.
// The test for video file is based on the file extension.
func (encoder *Encoder) getVideoFiles() ([]string, error) {
	var videoFiles []string
	err := filepath.Walk(encoder.Source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// TODO log error and return SkipDir?
			return errors.Wrap(err, "failed to walk source directory")
		}
		if encoder.isVideoFile(path, info) {
			videoFiles = append(videoFiles, path)
		}
		return nil
	})
	return videoFiles, err
}

// isVideoFile returns whether the given path is a video file based on the extension.
func (encoder *Encoder) isVideoFile(path string, info os.FileInfo) bool {
	if info.IsDir() {
		return false
	}

	ext := filepath.Ext(path)
	for _, ve := range encoder.config.SourceExtensions {
		if ve == ext {
			return true
		}
	}
	return false
}

// getDestinationPath returns the path to the new file based at the destination.
func (encoder *Encoder) getDestinationPath(videoPath string) (string, error) {
	rel, err := filepath.Rel(encoder.Source, videoPath)
	if err != nil {
		return "", errors.Wrap(err, "failed to find relative path to source file")
	}
	destinationPath := filepath.Join(encoder.Target, rel)
	ext := filepath.Ext(destinationPath)
	destinationPath = strings.TrimSuffix(destinationPath, ext)
	destinationPath = destinationPath + encoder.config.TargetExtension
	return destinationPath, nil
}
