package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gary-y-chang/concurrency/patterns/runner"
)

const timeout = 7 * time.Second

func main() {
	fmt.Println("Start Concurrency Pattern TaskRunner ....")
	fmt.Printf("Select a pattern to run.\n '%s' for TaskRunner\n '%s' for Pool\n", "R", "P")

	reader := bufio.NewReader(os.Stdin)
	char, _, err := reader.ReadRune()
	if err != nil {
		fmt.Println(err)
	}
	
    switch char {
	case 'R':
    	fmt.Printf("'%s' pressed.  Start Pattern TaskRunner ....", string(char))
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
	
case 'P':
		fmt.Printf("'%s' pressed.  Start Pattern Pool ....", string(char))
	}


    
}

func createTask() func(int) {
	log.Printf("---- > Task created.")
	return func(id int) {
		log.Printf("Task_ID:#%d getting started.", id)
		time.Sleep(time.Duration(2) * time.Second)
		log.Printf("Task_ID:#%d finished.", id)
	}
}