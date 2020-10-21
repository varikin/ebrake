package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

type Encoder struct {
	config *Config
}

type Video struct {
	source string
	target string
}

// EncodeFiles encodes all the video files in the source directory and places them
// in the destination directory.
func (encoder *Encoder) EncodeFiles() error {

	info, err := os.Stat(encoder.config.Target)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return errors.New(encoder.config.Target + " is not a directory")
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
		videos = append(videos, Video{
			source: videoFile,
			target: targetFile,
		})
	}
	if len(videos) == 0 {
		fmt.Println("Did not find any videos re-encode.")
	} else {
		for _, video := range videos {
			fmt.Println(encoder.getEncodeCommand(video))
		}
	}

	return nil
}

func (encoder *Encoder) getEncodeCommand(video Video) string {
	return fmt.Sprintf("%s %s -i %s -o %s",
		encoder.config.HandBrakeCommand,
		encoder.config.HandBrakeOptions,
		video.source,
		video.target,
	)
}

// getVideoFiles finds and returns the list of video files in the given directory.
// The test for video file is based on the file extension.
func (encoder *Encoder) getVideoFiles() ([]string, error) {
	var videoFiles []string
	err := filepath.Walk(encoder.config.Source, func(path string, info os.FileInfo, err error) error {
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
	rel, err := filepath.Rel(encoder.config.Source, videoPath)
	if err != nil {
		return "", errors.Wrap(err, "failed to find relative path to source file")
	}
	destinationPath := filepath.Join(encoder.config.Target, rel)
	ext := filepath.Ext(destinationPath)
	destinationPath = strings.TrimSuffix(destinationPath, ext)
	destinationPath = destinationPath + encoder.config.TargetExtension
	return destinationPath, nil
}
