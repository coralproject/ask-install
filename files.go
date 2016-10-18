//go:generate go-bindata -pkg $GOPACKAGE -o templates.go templates/
package main

import (
	"fmt"
	"html/template"
	"os"
)

const (
	configFilename      = "askenv"
	setupScriptFilename = "setup.sh"
)

var (
	dockerComposeTemplate = template.Must(template.New("docker-compose.yml").Parse(string(MustAsset("templates/docker-compose.yml"))))
	setupScriptTemplate   = template.Must(template.New(setupScriptFilename).Parse(string(MustAsset("templates/" + setupScriptFilename))))
	caddyfileTemplate     = template.Must(template.New("Caddyfile").Parse(string(MustAsset("templates/Caddyfile"))))
)

// CreateDockerComposeFile will template out the docker-compose.yml file.
func CreateDockerComposeFile(cayport string) error {
	f, err := os.Create("docker-compose.yml")
	if err != nil {
		return err
	}
	defer f.Close()

	ctx := map[string]string{
		"Port":           cayport,
		"ConfigFilename": configFilename,
	}

	if err := dockerComposeTemplate.Execute(f, ctx); err != nil {
		return err
	}

	return nil
}

// CreateConfigFile will create the askenv configuration file.
func CreateConfigFile(config [][]string) error {
	f, err := os.Create(configFilename)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, configline := range config {
		if _, err := fmt.Fprintf(f, "%s=%s\n", configline[0], configline[1]); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprint(f, "\n"); err != nil {
		return err
	}

	return nil
}

// CreateSetupScript will create the setup script to start the application for
// the first time.
func CreateSetupScript(ctx map[string]string) error {
	f, err := os.Create(setupScriptFilename)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := setupScriptTemplate.Execute(f, ctx); err != nil {
		return err
	}

	return nil
}

// CreateCaddyFile will tempalte out the Caddyfile to be used by Caddy.
func CreateCaddyFile(hostname string) error {
	f, err := os.Create("Caddyfile")
	if err != nil {
		return err
	}
	defer f.Close()

	ctx := map[string]interface{}{
		"Hostname": hostname,
	}

	if err := caddyfileTemplate.Execute(f, ctx); err != nil {
		return err
	}

	return nil
}
