package main

import (
	"fmt"
	"log"
	"os"

	event "github.com/nogurenn/assorted-programs/simple-event-worker"
)

func main() {
	// Dependencies --------------------------------------
	logger := log.Default()
	logger.SetPrefix("[main] ")

	eventService := event.NewService()

	// Execute --------------------------------------------

	cwd, err := os.Getwd()
	if err != nil {
		logger.Println(err)
		return
	}

	// Simulate reading the events from anywhere that we could read from. We only need an io.Reader instance.
	// In this case, we're reading from a file.
	file, err := os.Open(fmt.Sprintf("%s/events.json", cwd))
	if err != nil {
		logger.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	events, err := eventService.ParseEvents(file)
	if err != nil {
		logger.Println(err)
		os.Exit(1)
	}

	// Normally, all the processing is already done in the eventService.ProcessEvents method,
	// and we don't need the events/accounts anymore.
	// Maybe it saves the results to a database or sends them to another service, or saves them to a file or whatever.
	// We're just printing the results here for demonstration.
	accounts, err := eventService.ProcessEvents(events)
	if err != nil {
		logger.Println(err)
		os.Exit(1)
	}

	for id, account := range accounts {
		logger.Printf("%s: {Status: %s, Balance: %d}\n", id, account.Status(), account.Balance())
	}

	os.Exit(0)
}
