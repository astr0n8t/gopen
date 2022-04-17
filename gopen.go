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

	"github.com/astr0n8t/gopen/definitions"
	"github.com/astr0n8t/gopen/modules"
	"github.com/astr0n8t/gopen/utilities"
)

func main() {

	confOptions := readConfig()

	results := utilities.InitHosts(confOptions.Variables)

	for command, options := range confOptions.Workflow {
		step := modules.GetModule(command, confOptions.Variables, options, results)
		results := step.RunModule()
		if results.Success {
			fmt.Println(step.GetOutput())
		} else {
			fmt.Println("Step " + command + " failed!")
		}
	}
}

// Processes the configuration file and command line arguments using Viper and PFlags
func readConfig() definitions.Config {

	// Set defaults
	viper.SetDefault("variables.addresses", "")
	viper.SetDefault("variables.root", false)

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
	pflag.String("variables.addresses", "", "Address(es) to scan")
	pflag.String("variables.ports", "", "Ports) to scan")
	pflag.Bool("variables.root", false, "Whether it should be ran as root")
	pflag.Parse()

	// Add the command line arguments to viper
	viper.BindPFlags(pflag.CommandLine)

	// Unmarshall the config file into the config struct
	var processedConfig definitions.Config
	err = viper.Unmarshal(&processedConfig)
	if err != nil {
		panic(fmt.Errorf("unable to unmarshall config file or command line arguments"))
	}

	// Check for required configuration options
	if processedConfig.Variables.Addresses == "" {
		panic(fmt.Errorf("one or more required arguments not supplied or config file could not be read\n required arguments: address"))
	}

	// Return the config struct with config data
	return processedConfig
}
