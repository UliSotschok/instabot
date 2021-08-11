package bot

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/ahmdrz/goinsta/v2"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
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
	flag.BoolVar(&run, "run", true, "Use this option to follow, like and comment")
	flag.BoolVar(&unfollow, "sync", false, "Use this option to unfollow those who are not following back")
	flag.BoolVar(&nomail, "nomail", true, "Use this option to disable the email notifications")
	flag.BoolVar(&dev, "dev", false, "Use this option to use the script in development mode : nothing will be done for real")
	flag.BoolVar(&logs, "logs", false, "Use this option to enable the logfile")
	flag.BoolVar(&noduplicate, "noduplicate", true, "Use this option to skip following, liking and commenting same user in this session")

	flag.Parse()

	// -logs enables the log file
	if logs {
		// Opens a log file
		t := time.Now()
		logFile, err := os.OpenFile("instabot-"+t.Format("2006-01-02-15-04-05")+".log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		check(err)
		defer func(logFile *os.File) {
			err := logFile.Close()
			if err != nil {
				log.Printf("Error closing file. error='%v'", err)
			}
		}(logFile)

		// Duplicates the writer to stdout and logFile
		mw := io.MultiWriter(os.Stdout, logFile)
		log.SetOutput(mw)
	}
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

func getInput(text string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf(text)
	input, err := reader.ReadString('\n')
	check(err)
	return strings.TrimSpace(input)
}

// Checks if the user is in the slice
func containsUser(slice []goinsta.User, user goinsta.User) bool {
	for _, currentUser := range slice {
		if currentUser.Username == user.Username {
			return true
		}
	}
	return false
}

func getInputf(format string, args ...interface{}) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf(format, args...)
	input, err := reader.ReadString('\n')
	check(err)
	return strings.TrimSpace(input)
}

// Same, with strings
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
