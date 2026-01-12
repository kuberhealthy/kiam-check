package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kuberhealthy/kuberhealthy/v3/pkg/checkclient"
	log "github.com/sirupsen/logrus"
)

// main loads configuration and executes the KIAM Lambda check.
func main() {
	// Parse configuration from environment variables.
	cfg, err := parseConfig()
	if err != nil {
		reportFailureAndExit(err)
		return
	}

	// Apply debug settings after parsing.
	applyDebugSettings(cfg)

	// Give the k8s API enough time to allocate IPs.
	time.Sleep(15 * time.Second)

	// Create an AWS session.
	sess, err := createAWSSession()
	if err != nil {
		reportFailureAndExit(fmt.Errorf("failed to create AWS session: %w", err))
		return
	}
	if sess == nil {
		reportFailureAndExit(fmt.Errorf("nil AWS session"))
		return
	}

	// Start listening for interrupts.
	signalChan := make(chan os.Signal, 2)
	go listenForInterrupts(signalChan)

	// Catch panics to report failures.
	defer handlePanic()

	// Create a Lambda client.
	lambdaClient := createLambdaClient(sess, cfg.AWSRegion)

	// Run the Lambda list check.
	err = runLambdaCheck(cfg, lambdaClient)
	if err != nil {
		reportFailureAndExit(fmt.Errorf("error occurred during Lambda check: %w", err))
		return
	}
	log.Infoln("AWS Lambda check successful.")

	// Report success to Kuberhealthy.
	reportSuccessAndExit()
}

// listenForInterrupts waits for an OS signal and exits when received.
func listenForInterrupts(signalChan chan os.Signal) {
	// Relay incoming OS interrupt signals to the signal channel.
	signal.Notify(signalChan, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT)
	sig := <-signalChan
	log.Infoln("Received an interrupt signal from the signal channel.")
	log.Debugln("Signal received was:", sig.String())

	// Clean up pods here.
	log.Infoln("Shutting down.")

	os.Exit(0)
}

// handlePanic reports panics as failures to Kuberhealthy.
func handlePanic() {
	// Recover from panics during the check run.
	recovered := recover()
	if recovered == nil {
		return
	}

	// Report the panic to Kuberhealthy.
	log.Infoln("Recovered panic:", recovered)
	reportFailureAndExit(fmt.Errorf("panic: %v", recovered))
}

// reportSuccessAndExit reports success to Kuberhealthy and exits.
func reportSuccessAndExit() {
	// Report the success to Kuberhealthy.
	err := checkclient.ReportSuccess()
	if err != nil {
		log.Fatalln("error reporting to kuberhealthy:", err.Error())
	}

	// Exit after reporting success.
	os.Exit(0)
}

// reportFailureAndExit reports failure to Kuberhealthy and exits.
func reportFailureAndExit(err error) {
	// Log the error and report failure.
	log.Errorln("Reporting errors to Kuberhealthy:", err.Error())
	reportErr := checkclient.ReportFailure([]string{err.Error()})
	if reportErr != nil {
		log.Fatalln("error reporting to kuberhealthy:", reportErr.Error())
	}

	// Exit after reporting failure.
	os.Exit(0)
}
