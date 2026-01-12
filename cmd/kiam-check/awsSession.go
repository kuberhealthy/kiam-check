package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	log "github.com/sirupsen/logrus"
)

// createAWSSession creates and returns an AWS session.
func createAWSSession() (*session.Session, error) {
	// Build an AWS session with verbose credential errors.
	log.Infoln("Building AWS session.")
	return session.NewSession(aws.NewConfig().WithCredentialsChainVerboseErrors(true))
}
