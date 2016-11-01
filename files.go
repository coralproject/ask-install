//go:generate go-bindata -pkg $GOPACKAGE -o templates.go templates/
package main

import (
	"encoding/json"
	"os"
	"text/template"
)

const (
	askInstallerStateFilename = "ask-install.json"
	setupScriptFilename       = "setup.sh"
)

var (
	dockerComposeTemplate = template.Must(template.New("docker-compose.yml").Parse(string(MustAsset("templates/docker-compose.yml"))))
	setupScriptTemplate   = template.Must(template.New(setupScriptFilename).Parse(string(MustAsset("templates/" + setupScriptFilename))))
	caddyfileTemplate     = template.Must(template.New("Caddyfile").Parse(string(MustAsset("templates/Caddyfile"))))
)

// CreateDockerComposeFile will template out the docker-compose.yml file.
func CreateDockerComposeFile(config Config) error {
	f, err := os.Create("docker-compose.yml")
	if err != nil {
		return err
	}
	defer f.Close()

	if err := dockerComposeTemplate.Execute(f, config); err != nil {
		return err
	}

	return nil
}

// CreateSetupScript will create the setup script to start the application for
// the first time.
func CreateSetupScript(config Config) error {
	f, err := os.Create(setupScriptFilename)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := setupScriptTemplate.Execute(f, config); err != nil {
		return err
	}

	return nil
}

// CreateCaddyFile will tempalte out the Caddyfile to be used by Caddy.
func CreateCaddyFile(config Config) error {
	f, err := os.Create("Caddyfile")
	if err != nil {
		return err
	}
	defer f.Close()

	if err := caddyfileTemplate.Execute(f, config); err != nil {
		return err
	}

	return nil
}

// LoadAskInstallState loads the state of the configuration from the filesystem.
func LoadAskInstallState() (*Config, error) {
	f, err := os.Open(askInstallerStateFilename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var config Config
	if err := json.NewDecoder(f).Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// CreateAskInstallState saves the state of the configuration to the filesystem.
func CreateAskInstallState(config Config) error {
	f, err := os.Create(askInstallerStateFilename)
	if err != nil {
		return err
	}
	defer f.Close()

	e := json.NewEncoder(f)

	e.SetIndent("", "  ")

	if err := e.Encode(config); err != nil {
		return err
	}

	return nil
}
