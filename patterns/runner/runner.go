package runner

import (
	"errors"
	"os"
	"os/signal"
	"time"
)

type TaskRunner struct {
	
	interrupt chan os.Signal

	complete chan error

	timeout <- chan time.Time
	
	tasks []func(int)
}

var ErrTimeout = errors.New("received timeout")
var ErrInterrupt = errors.New("received interrupt")

func (runner *TaskRunner) gotInterrupt() bool {
	select{
	case <- runner.interrupt:
		signal.Stop(runner.interrupt)
		return true

	default:
		return false	
	}
}

func (runner *TaskRunner) run() error {
	for id, task := range runner.tasks {
		if runner.gotInterrupt() {
			return ErrInterrupt
		}
	task(id)
}
	return nil
}

func New(d time.Duration) *TaskRunner {
	return &TaskRunner{
		interrupt: make(chan os.Signal, 1),
		complete: make(chan error),
		timeout: time.After(d),
	}
}

func (runner *TaskRunner) Add(tasks ...func(int)) {
	runner.tasks = append(runner.tasks, tasks...)
}

func (runner *TaskRunner) Start() error {
	signal.Notify(runner.interrupt, os.Interrupt)

	go func () {
		runner.complete <-runner.run()
	}()

	select{
		case err := <-runner.complete:
			return err
    	
		case <-runner.timeout:
			return ErrTimeout
	}
}

