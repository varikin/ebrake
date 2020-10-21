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
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	HandBrakeCommand string
	HandBrakeOptions string
	Source           string
	Target           string
	SourceExtensions []string
	TargetExtension  string
}

var cfgFile string

// init initialized cobra for the CLI flags
func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(
		&cfgFile,
		"config",
		"",
		"config file (default is $HOME/.ebrake.yaml)",
	)
}

// main let's go dude
func main() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetDefault("SourceExtensions", []string{".mp4", ".mkv", ".avi"})
	viper.SetDefault("TargetExtension", ".mp4")
	viper.SetDefault("HandBrakeCommand", "HandBrakeCLI.exe")
	viper.SetDefault("HandBrakeOptions", "--encoder x264 --encoder-preset fast --optimize")

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".ebrake" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".ebrake")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

// rootCmd represents the base CLI command
var rootCmd = &cobra.Command{
	Use:   "ebrake",
	Short: "Re-encodes a directory of movie files using Handbrake.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		config := Config{
			Source: args[0],
			Target: args[1],
		}
		if err := viper.Unmarshal(&config); err != nil {
			fmt.Println("Error: " + err.Error())
			os.Exit(1)
		}

		encoder := Encoder{config: &config}
		if err := encoder.EncodeFiles(); err != nil {
			fmt.Println("Error: " + err.Error())
			os.Exit(1)
		}
	},
}
