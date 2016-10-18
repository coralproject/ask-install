package main

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/fatih/color"
)

func main() {
	fmt.Println("Coral Project ASK Installer")

	color.Cyan("\nGeneral Configuration\n")

	var config struct {
		Hostname          string
		RootURL           string
		S3Bucket          string
		AWSRegion         string
		AWSAccessKeyID    string
		AWSAccessKey      string
		RecaptchaSecret   string
		GoogleAnalyticsID string
		Email             string
		Password          string
		DisplayName       string
		Port              string
	}

	fmt.Println(`
This is where you can specify the host on which the provided server will bind
to. If you specify the host with a port, it will specifically bind to that port,
otherwise, port 80, 443 will be bound to
`)

	for {
		config.Hostname = StringRequired("What's the external hostname of this machine?")

		if strings.Contains(config.Hostname, "http") {
			color.Red("Hostname can't contain the scheme (http://, https://)")
			continue
		}

		if strings.Contains(config.Hostname, ":") {
			_, port, err := net.SplitHostPort(config.Hostname)
			if err != nil {
				port = "80"
			}

			config.Port = port

		} else {
			config.Port = "80"
		}

		sslEnabled := Confirm("Do you want SSL enabled?")

		if sslEnabled {
			config.RootURL = "https://" + config.Hostname
		} else {
			config.RootURL = "http://" + config.Hostname
		}

		if ok := Confirm("External URL will be \"%s\", is that ok?", config.RootURL); ok {
			break
		}
	}

	if ok := Confirm("Do you want to enable recaptcha?"); ok {
		config.RecaptchaSecret = StringRequired("What is the recaptcha server secret?")
	}

	if ok := Confirm("Do you want to enable Google Analytics?"); ok {
		config.GoogleAnalyticsID = StringRequired("What is the Google Analytics ID?")
	}

	// TODO: PROJECT
	// TODO: ENV

	color.Cyan("\nAmazon\n")

	if ok := Confirm("Do you want forms uploaded to S3?"); ok {
		config.S3Bucket = StringRequired("What's the S3 Bucket we can upload forms?")
		config.AWSRegion = StringRequired("What's the S3 Region for this bucket?")
		config.AWSAccessKeyID = StringRequired("What's the AWS_ACCESS_KEY_ID with write access?")
		config.AWSAccessKey = StringRequired("What's the AWS_ACCESS_KEY associated with this AWS_ACCESS_KEY_ID?")
	}

	color.Cyan("\nAuth\n")

	config.DisplayName = StringRequired("What's the name for the user account?")
	config.Email = StringRequired("What's the email address for the user account?")
	config.Password = PasswordMasked("What's the password for the account?")

	privateKey, publicKey, err := GenerateKeys()
	if err != nil {
		color.Red("Couldn't create keys: %s", err.Error())
		os.Exit(1)
	}

	sessionSecret, _, err := GenerateKeys()
	if err != nil {
		color.Red("Couldn't create keys: %s", err.Error())
		os.Exit(1)
	}

	color.Cyan("\nCreating Files\n")

	setupContext := map[string]string{
		"Email":       config.Email,
		"Password":    config.Password,
		"DisplayName": config.DisplayName,
	}

	if err := CreateSetupScript(setupContext); err != nil {
		color.Red("Couldn't create %s: %s", setupScriptFilename, err.Error())
		os.Exit(1)
	}
	color.Green("Created: %s", setupScriptFilename)

	if err := CreateCaddyFile(config.Hostname); err != nil {
		color.Red("Couldn't create Caddyfile: %s", err.Error())
		os.Exit(1)
	}
	color.Green("Created: Caddyfile")

	configMap := [][]string{
		// Coral Auth
		{"CORAL_AUTH_PRIVATE_KEY", privateKey},
		{"CORAL_AUTH_PUBLIC_KEY", publicKey},
		{"CORAL_AUTH_SESSION_SECRET", sessionSecret},
		{"CORAL_AUTH_ALLOWED_CLIENTS", "cay " + config.RootURL + "/callback"},
		{"CORAL_AUTH_ROOT_URL", config.RootURL + "/auth"},

		// askd
		{"ASK_AUTH_PUBLIC_KEY", publicKey},
		{"ASK_RECAPTCHA_SECRET", ""},

		// Cay
		{"GAID", config.GoogleAnalyticsID},
		{"AUTH_CLIENT_ID", "cay"},
		{"AUTH_AUTHORITY", config.RootURL + "/auth/connect"},
		{"ELKHORN_URL", config.RootURL + "/elkhorn"},
		{"ASK_URL", config.RootURL + "/askd"},

		// Elkhorn
		{"S3_BUCKET", config.S3Bucket},
		{"AWS_REGION", config.S3Bucket},
		{"AWS_ACCESS_KEY_ID", config.S3Bucket},
		{"AWS_ACCESS_KEY", config.S3Bucket},
		{"RECAPTCHA", config.RecaptchaSecret},
	}

	if err := CreateConfigFile(configMap); err != nil {
		color.Red("Couldn't create %s: %s", configFilename, err.Error())
		os.Exit(1)
	}
	color.Green("Created: %s", configFilename)

	if err := CreateDockerComposeFile(config.Port); err != nil {
		color.Red("Couldn't create docker-compose.yml: %s", err.Error())
		os.Exit(1)
	}
	color.Green("Created: docker-compose.yml")

	color.Green("\n\nFinished! Run the following to start using Ask!\n\n\tbash %s\n", setupScriptFilename)
}
