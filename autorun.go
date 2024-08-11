package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func startAutorun(drivePath string) {
	// Check if the drive path exists
	if _, err := os.Stat(drivePath); os.IsNotExist(err) {
		fmt.Printf("Drive path %s does not exist\n", drivePath)
		return
	}

	// Read the config file
	conf, err := setupConfig(drivePath + "/.autorun.toml")
	if err != nil {
		fmt.Printf("Error reading config file: %s\n", err)
		return
	}

	if conf.Autorun == "" {
		fmt.Println("No autorun program specified")
		return
	}

	// Set the environment variables
	conf = setupEnvironment(conf)
	// convert the environment map to a slice
	env := make([]string, 0, len(conf.Environment))

	// start building the command
	cmd := exec.Command(conf.Autorun)

	if conf.Isolate {
		cmd.Env = env
	} else {
		cmd.Env = append(os.Environ(), env...)
	}

	// Set the working directory
	if conf.WorkDir != "" {
		cmd.Dir = conf.WorkDir
	} else {
		cmd.Dir = drivePath
	}

	if !filepath.IsAbs(conf.Autorun) {
		conf.Autorun = filepath.Join(drivePath, conf.Autorun)
	}

	// Start the autorun program
	err = cmd.Start()
	if err != nil {
		fmt.Printf("Error starting autorun program: %s\n", err)
		return
	}

}
