package bot

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"time"
)

// check will log.Fatal if err is an error
func check(err error) {
	if err != nil {
		log.Fatal("ERROR:", err)
	}
}

// Parses the options given to the script
func parseOptions() {
	flag.BoolVar(&dev, "dev", false, "Use this option to use the script in development mode : nothing will be done for real")
	flag.Parse()
}

func setupLogging() {
	// Opens a log file
	t := time.Now()
	logFile, err := os.OpenFile("instabot-"+t.Format("2006_01_02")+".log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	check(err)

	// Duplicates the writer to stdout and logFile
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
}

// Retries the same function [function], a certain number of times (maxAttempts).
// It is exponential : the 1st time it will be (sleep), the 2nd time, (sleep) x 2, the 3rd time, (sleep) x 3, etc.
// If this function fails to recover after an error, it will send an email to the address in the config file.
func retry(maxAttempts int, sleep time.Duration, function func() error) (err error) {
	for currentAttempt := 0; currentAttempt < maxAttempts; currentAttempt++ {
		err = function()
		if err == nil {
			return
		}
		for i := 0; i <= currentAttempt; i++ {
			time.Sleep(sleep)
		}
		log.Println("Retrying after error:", err)
	}

	return fmt.Errorf("After %d attempts, last error: %s", maxAttempts, err)
}

func containsString(slice []string, user string) bool {
	for _, currentUser := range slice {
		if currentUser == user {
			return true
		}
	}
	return false
}

func doPauseAfterAction(config scheduling) {
	duration := time.Duration(rand.Intn(config.PauseAfterActionInS)+config.PauseAfterActionInS) * time.Second
	log.Printf("Sleeping for: %v\n", duration)
	time.Sleep(duration)
}

func doPauseAfterBatch(config scheduling) {
	duration := time.Duration(rand.Intn(config.PauseAfterBatchInM)+config.PauseAfterBatchInM) * time.Minute
	log.Printf("Sleeping for: %v\n", duration)
	time.Sleep(duration)
}
