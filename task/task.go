package task

import (
	"fmt"
	"go-thread-monitoring/sender"
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
		taskChannel: make(chan struct{}, 5),
		connManager: connManager,
	}
}

func (tm *TaskManager) StartRandomTask() {
	tm.taskChannel <- struct{}{} // push into channel random elements

	// lock the resource, update count of process in "PROCESSING"
	tm.mutex.Lock()
	tm.activeTask++
	tm.mutex.Unlock()

	// simulate processing
	duration := time.Duration(rand.Intn(5)+1) * time.Second
	time.Sleep(duration)

	// lock the resource, since task complete its process, remove 1 from process count
	tm.mutex.Lock()
	tm.activeTask--
	tm.mutex.Unlock()

	// send task to esit queue
	message := fmt.Sprintf("Task completed in %v seconds ", duration.Seconds())
	tm.connManager.Publish("task", message)
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
