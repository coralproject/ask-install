package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/fatih/color"
)

// Config stores the answers to the questions made from the interactive console.
type Config struct {
	Hostname          string
	RootURL           string
	UseS3             bool
	UseSSL            bool
	S3Bucket          string
	AWSRegion         string
	AWSAccessKeyID    string
	AWSAccessKey      string
	RecaptchaSecret   string
	GoogleAnalyticsID string
	Email             string `json:"-"`
	Password          string `json:"-"`
	DisplayName       string `json:"-"`
	Port              string
	AuthPublicKey     string
	AuthPrivateKey    string
	SessionSecret     string
	Channel           string
}

// GetConfigurationFromInteractive uses prompts to request the configuration
// options.
func GetConfigurationFromInteractive() (*Config, error) {
	color.Cyan("\nGeneral Configuration\n")

	var config Config

	if useStable := Confirm("Do you want to use the stable version of ask?"); useStable {
		config.Channel = "release"
	} else {
		config.Channel = "latest"
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

		config.UseSSL = Confirm("Do you want SSL enabled?")

		if config.UseSSL {
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

	color.Cyan("\nAmazon\n")

	if config.UseS3 = Confirm("Do you want forms uploaded to S3?"); config.UseS3 {
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
		return nil, fmt.Errorf("Couldn't create keys: %s", err.Error())
	}

	config.AuthPrivateKey = privateKey
	config.AuthPublicKey = publicKey

	sessionSecret, _, err := GenerateKeys()
	if err != nil {
		return nil, fmt.Errorf("Couldn't create keys: %s", err.Error())
	}

	config.SessionSecret = sessionSecret

	return &config, nil
}
