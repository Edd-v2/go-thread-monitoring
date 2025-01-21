package task

import (
	"fmt"
	"go-thread-monitoring/sender"
	"log"
	"math/rand"
	"sync"
	"time"
)

type TaskManager struct {
	mutex       sync.Mutex
	activeTask  int
	taskChannel chan struct{}
	connManager *sender.ConnectionManager
}

func NewTaskManager(connManager *sender.ConnectionManager) *TaskManager {
	return &TaskManager{
		taskChannel: make(chan struct{}, 3),
		connManager: connManager,
	}
}

func (tm *TaskManager) StartRandomTask() {
	log.Println("Attempting to start a new task...")

	select {
	case tm.taskChannel <- struct{}{}: //if push into channel of random elements is possible
		log.Println("Task added to the channel")

		// lock the resource, update count of process in "PROCESSING"
		tm.mutex.Lock()
		tm.activeTask++
		tm.mutex.Unlock()
		log.Println("Active task count incremented")

		// simulate processing
		duration := time.Duration(rand.Intn(15)+1) * time.Second
		log.Printf("Stimated Task duration for %v seconds...\n", duration.Seconds())
		time.Sleep(duration)

		<-tm.taskChannel //remove task from channell since processing is terminated

		// lock the resource, since task complete its process, remove 1 from process count
		tm.mutex.Lock()
		tm.activeTask--
		tm.mutex.Unlock()
		log.Println("Task removed - count decremented")

		// send task to esit queue
		message := fmt.Sprintf("Task completed in %v seconds ", duration.Seconds())
		tm.connManager.Publish("task", message)
	default:
		log.Println("Unable to start task: channel is full. Try again later.") // or create another channel
	}

}

func (tm *TaskManager) GetTaskStatus() int {
	tm.mutex.Lock()
	defer tm.mutex.Unlock() // ensure unlock after process return
	return tm.activeTask
}

func (tm *TaskManager) ClearAllChannel() {

	/**
		if there is not concurrency

		for len(ch) > 0 {
	  		<-ch
		}

	*/

L:
	for {
		select {
		case _, ok := <-tm.taskChannel: // read from channel
			if !ok { //if closed immediately return err
				break L
			}
		default: // nothing in ch for now
			break L
		}
	}

}
