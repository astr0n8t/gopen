/*
----------------goPen----------------
A YAML defined penetration testing
workflow tool written in Golang
--------------------------------------------
Written by Nathan Higley (@astr0n8t)
nathan@nathanhigley.com
https://nathanhigley.com
https://github.com/astr0n8t
--------------------------------------------
*/

package main

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {

	confOptions := readConfig()

	fmt.Println(confOptions.Address)
}

// A struct to store configuration options
type config struct {
	Executable string
	Address    string
	Flags      string
	Root       bool
}

// Processes the configuration file and command line arguments using Viper and PFlags
func readConfig() config {

	// Set defaults
	viper.SetDefault("Executable", "nmap")
	viper.SetDefault("Address", "")
	viper.SetDefault("Flags", "")
	viper.SetDefault("Root", false)

	// Set default config directories
	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/.config/gopen")
	viper.AddConfigPath("/etc/gopen/")
	viper.AddConfigPath(".")

	// Read the config
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Printf("Config file not found, reading command line arguments only.\n")
		} else {
			panic(fmt.Errorf("fatal error config file: %s", err))
		}
	}

	// Add and read command line arguments
	pflag.String("address", "", "Address(es) to scan")
	pflag.String("executable", "", "Path to an executable to run")
	pflag.String("flags", "", "Flags to pass to the executable")
	pflag.Bool("root", false, "Whether it should be ran as root")
	pflag.Parse()

	// Add the command line arguments to viper
	viper.BindPFlags(pflag.CommandLine)

	// Unmarshall the config file into the config struct
	var processedConfig config
	err = viper.Unmarshal(&processedConfig)
	if err != nil {
		panic(fmt.Errorf("unable to unmarshall config file or command line arguments"))
	}

	// Check for required configuration options
	if processedConfig.Address == "" {
		panic(fmt.Errorf("one or more required arguments not supplied or config file could not be read\n required arguments: address"))
	}

	// Return the config struct with config data
	return processedConfig
}
