package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gary-y-chang/concurrency/patterns/pool"
	"github.com/gary-y-chang/concurrency/patterns/runner"
)

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
		timeout := 7 * time.Second
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
			const maxGoroutines = 10 
			const pooledResources = 2
			var wg sync.WaitGroup
			wg.Add(maxGoroutines)
			
			p, err := pool.New(createDbConnection, pooledResources)
			if err != nil {
				log.Println(err)
			}

			for query :=0 ; query < maxGoroutines ; query++ {
				go func(queryId int) {
                    conn, err := p.Acquire()
					if err != nil {
						log.Println(err)
						return
					}
					defer p.Release(conn)

					
					log.Printf("Query[%d] with DbConnection[%d] in process\n\n", queryId, conn.(*dbConnection).ID )		
							
					time.Sleep(time.Duration(rand.Intn(2000))*time.Millisecond)	
					wg.Done()
				}(query)

				time.Sleep(time.Duration(1 * time.Second))	
			}

           wg.Wait()
		  
		   p.Close()
		   log.Println("-----Close the Connection Pool.-----")
	}
	
}


type dbConnection struct {
	ID int32
}

func (dbConn *dbConnection) Close() error {
	fmt.Printf("Close dbConnection ID: %d.", dbConn.ID)
	return nil
}

var idCounter int32

func createDbConnection() (io.Closer, error) {
	id := atomic.AddInt32(&idCounter, 1)
	log.Printf("New DB Connection created with ID: %d.", id)
	return &dbConnection{id}, nil
}

func createTask() func(int) {
	log.Printf("---- > Task created.")
	return func(id int) {
		log.Printf("Task_ID:#%d getting started.", id)
		time.Sleep(time.Duration(2) * time.Second)
		log.Printf("Task_ID:#%d finished.", id)
	}
}