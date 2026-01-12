package main

import (
	"fmt"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
)

// CheckConfig stores configuration for the KIAM check.
type CheckConfig struct {
	// AWSRegion is the region to query for Lambda functions.
	AWSRegion string
	// ExpectedLambdaCount is an optional count to enforce.
	ExpectedLambdaCount int
	// Debug enables debug logging when true.
	Debug bool
}

// defaultAWSRegion is the default region when none is provided.
const defaultAWSRegion = "us-west-2"

// parseConfig parses environment variables into a CheckConfig.
func parseConfig() (*CheckConfig, error) {
	// Parse debug settings.
	debug, err := parseDebugSetting()
	if err != nil {
		return nil, err
	}

	// Parse the AWS region.
	awsRegion := defaultAWSRegion
	awsRegionEnv := os.Getenv("AWS_REGION")
	if len(awsRegionEnv) != 0 {
		awsRegion = awsRegionEnv
		log.Infoln("Parsed AWS_REGION:", awsRegion)
	}

	// Parse the expected Lambda count.
	expectedLambdaCount := 0
	lambdaCountEnv := os.Getenv("LAMBDA_COUNT")
	if len(lambdaCountEnv) != 0 {
		count, parseErr := strconv.Atoi(lambdaCountEnv)
		if parseErr != nil {
			return nil, fmt.Errorf("error occurred attempting to parse LAMBDA_COUNT: %w", parseErr)
		}
		expectedLambdaCount = count
		log.Infoln("Parsed LAMBDA_COUNT:", expectedLambdaCount)
	}

	// Assemble configuration.
	cfg := &CheckConfig{}
	cfg.AWSRegion = awsRegion
	cfg.ExpectedLambdaCount = expectedLambdaCount
	cfg.Debug = debug

	return cfg, nil
}

// parseDebugSetting parses the DEBUG environment variable.
func parseDebugSetting() (bool, error) {
	// Default to disabled debug logging.
	debug := false

	// Parse DEBUG when provided.
	debugEnv := os.Getenv("DEBUG")
	if len(debugEnv) != 0 {
		parsedDebug, err := strconv.ParseBool(debugEnv)
		if err != nil {
			return false, fmt.Errorf("failed to parse DEBUG environment variable: %w", err)
		}
		debug = parsedDebug
	}

	return debug, nil
}

// applyDebugSettings updates logrus based on the debug flag.
func applyDebugSettings(cfg *CheckConfig) {
	// Enable debug logging when requested.
	if !cfg.Debug {
		return
	}

	// Apply logrus debug settings.
	log.Infoln("Debug logging enabled.")
	log.SetLevel(log.DebugLevel)
	log.Debugln(os.Args)
}
