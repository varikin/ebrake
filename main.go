/*
MIT License

Copyright (c) 2020 John Shimek

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
package main

import (
	"flag"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

var Version = "dev"

type Config struct {
	HandBrakeCommand string   `yaml:"handBrakeCommand"`
	HandBrakeOptions string   `yaml:"handBrakeOptions"`
	SourceExtensions []string `yaml:"sourceExtensions"`
	TargetExtension  string   `yaml:"targetExtensions"`
}

// main let's go dude
func main() {
	helpPtr := flag.Bool("help", false, "prints this help message")
	versionPtr := flag.Bool("version", false, "version for ebrake")
	configFile := flag.String("config", "", "config file (default is $HOME/.ebrake.yaml")
	flag.Parse()

	if *helpPtr {
		printUsage()
		os.Exit(0)
	}

	if *versionPtr {
		fmt.Println("ebrake version " + Version)
		os.Exit(0)
	}


	args := flag.Args()
	if len(args) != 2 {
		printUsage()
		os.Exit(1)
	}

	config, err := loadConfig(*configFile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	encoder := Encoder{
		config: config,
		Source: args[0],
		Target: args[1],
	}
	if err = encoder.EncodeFiles(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Re-encodes a directory of movie files using Handbrake.")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("\tebrake [-config=config-file] <source directory> <target directory>")
	fmt.Println("\tebrake -help")
	fmt.Println("\tebrake -version")
	fmt.Println("")
	fmt.Println("Options:")
	flag.PrintDefaults()
}

func loadConfig(configFile string) (*Config, error) {
	// Create the config with defaults
	config := &Config{
		HandBrakeCommand: "HandBrakeCLI.exe",
		HandBrakeOptions: "--encoder x264 --encoder-preset fast --optimize",
		SourceExtensions: []string{".mp4", ".mkv", ".avi"},
		TargetExtension:  ".mp4",
	}

	// A filename was given so it's required
	if configFile != "" {
		err := readConfigFile(config, configFile)
		if err != nil {
			return nil, errors.Wrap(err, "Unable to read config file: " + configFile)
		}
		return config, nil
	}

	// Fallback to look for one in the home directory, but not required

	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		return nil, errors.Wrap(err, "Unable to find home directory")
	}
	configFile = filepath.Join(home, ".ebrake.yaml")
	err = readConfigFile(config, configFile)
	if err != nil && !os.IsNotExist(err){
		// Skipping File Does Not Exist error
		return nil, errors.Wrap(err, "Unable to read config file: " + configFile)
	}
	// Return the config, either with defaults or with values from a file
	return config, nil
}

func readConfigFile(config *Config, filename string) error {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(file, config)
}

