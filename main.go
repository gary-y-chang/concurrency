package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gary-y-chang/concurrency/patterns/runner"
)

func main() {
	fmt.Printf("Start %s", "Concurrency Pattern TaskRunner ....")

	timeout := 6 * time.Second
    
	r := runner.New(timeout)

	r.Add(createTask(), createTask(), createTask())
	
	if err := r.Start(); err != nil {
		switch err {
		case runner.ErrTimeout:
			log.Printf("Terminating due to timeout.")
			os.Exit(1)
		
		case runner.ErrInterrupt:
			log.Printf("Terminating due to interrupt.")
			os.Exit(2)
		}
	} 

	log.Printf("All Tasks completed.")
}

func createTask() func(int) {
	return func(id int) {
		log.Printf("Processing Task_ID:#%d", id)
		time.Sleep(time.Duration(2) * time.Second)
	}
}