package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	log "github.com/sirupsen/logrus"
)

// createLambdaClient creates and returns an AWS Lambda client.
func createLambdaClient(sess *session.Session, region string) *lambda.Lambda {
	// Build a Lambda client.
	log.Infoln("Building Lambda client for", region, "region.")
	return lambda.New(sess, &aws.Config{Region: aws.String(region)})
}

// listLambdas returns a list of Lambda functions.
func listLambdas(client *lambda.Lambda) ([]*lambda.FunctionConfiguration, error) {
	// Create an output object for the list query.
	functions := make([]*lambda.FunctionConfiguration, 0)

	// Query Lambdas.
	log.Infoln("Querying Lambda functions.")
	results, err := client.ListFunctions(&lambda.ListFunctionsInput{
		MaxItems: aws.Int64(100),
	})
	if err != nil {
		return functions, err
	}

	// Append initial results.
	log.Infoln("Queried", len(results.Functions), "Lambdas.")
	functions = append(functions, results.Functions...)

	// Keep querying until the API gives us everything.
	for results.NextMarker != nil {
		log.Debugln("There are more results to be queried.")

		marker := *results.NextMarker
		results, err = client.ListFunctions(&lambda.ListFunctionsInput{
			MaxItems: aws.Int64(100),
			Marker:   aws.String(marker),
		})
		if err != nil {
			return functions, err
		}

		// Append subsequent results.
		log.Infoln("Queried", len(results.Functions), "Lambdas.")
		functions = append(functions, results.Functions...)
	}

	return functions, nil
}

// runLambdaCheck performs the Lambda list check and returns an error on failure.
func runLambdaCheck(cfg *CheckConfig, client *lambda.Lambda) error {
	// List Lambda functions.
	functions, err := listLambdas(client)
	if err != nil {
		return err
	}
	log.Infoln("Found", len(functions), "Lambdas.")

	// Validate against the expected count when configured.
	if cfg.ExpectedLambdaCount != 0 {
		if len(functions) != cfg.ExpectedLambdaCount {
			return fmt.Errorf("mismatching count of Lambdas -- expected %d, but got %d", cfg.ExpectedLambdaCount, len(functions))
		}
		return nil
	}

	// Pass when any Lambda is found.
	if len(functions) != 0 {
		return nil
	}

	// Fail when no Lambdas are found.
	return fmt.Errorf("could not find any Lambdas")
}
