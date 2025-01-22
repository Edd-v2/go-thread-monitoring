package main

import (
	"go-thread-monitoring/sender"
	"go-thread-monitoring/task"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	nats_url := "nats://demo.nats.io:4222"

	connManager := sender.NewConnectionManager()
	err := connManager.Connect(nats_url)
	if err != nil {
		log.Fatal("Could not connect to nats server")
	}

	init_api_and_processor(connManager)
}
func init_api_and_processor(connManager *sender.ConnectionManager) {
	r := gin.Default()
	task_manager := task.NewTaskManager(connManager)

	r.POST("/start-task", func(c *gin.Context) {
		go task_manager.StartRandomTask()
		c.JSON(http.StatusOK, gin.H{"status": "Task started"})
	})

	r.GET("/clear", func(c *gin.Context) {
		task_manager.ClearAllChannel()
		c.JSON(http.StatusOK, gin.H{"esit": "All channel process are empty"})
	})

	r.GET("/status", func(c *gin.Context) {
		status := task_manager.GetTaskStatus()
		c.JSON(http.StatusOK, gin.H{"active_tasks": status})
	})

	r.Run(":8080")

}
