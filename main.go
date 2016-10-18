package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fatih/color"
)

func main() {
	var useState = flag.Bool("s", false, "Use the state file in the current directory")

	flag.Parse()

	fmt.Println("Coral Project ASK Installer")

	var config *Config
	var err error

	if *useState {
		config, err = LoadAskInstallState()
		if err != nil {
			color.Red("Couldn't load state file: %s", err.Error())
			os.Exit(1)
		}

	} else {
		config, err = GetConfigurationFromInteractive()
		if err != nil {
			color.Red(err.Error())
			os.Exit(1)
		}
	}

	color.Cyan("\nCreating Files\n")
	if err := CreateFiles(*config); err != nil {
		color.Red("%s", err.Error())
		os.Exit(1)
	}

	color.Green("\n\nFinished! Run the following to start using Ask!\n\n\tbash %s\n", setupScriptFilename)

}

// CreateFiles creates all the files that are templated.
func CreateFiles(config Config) error {
	if err := CreateSetupScript(config); err != nil {
		return fmt.Errorf("Couldn't create %s: %s", setupScriptFilename, err.Error())
	}
	color.Green("Created: %s", setupScriptFilename)

	if err := CreateCaddyFile(config); err != nil {
		return fmt.Errorf("Couldn't create Caddyfile: %s", err.Error())
	}
	color.Green("Created: Caddyfile")

	if err := CreateDockerComposeFile(config); err != nil {
		return fmt.Errorf("Couldn't create docker-compose.yml: %s", err.Error())
	}
	color.Green("Created: docker-compose.yml")

	if err := CreateAskInstallState(config); err != nil {
		return fmt.Errorf("Couldn't create %s: %s", askInstallerStateFilename, err.Error())
	}
	color.Green("Created: %s", askInstallerStateFilename)

	return nil
}
