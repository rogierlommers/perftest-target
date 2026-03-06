package tasks

import (
	"fmt"
	"math/rand"

	"github.com/gin-gonic/gin"
)

func GETTasks(ctx *gin.Context) {
	task := generateRandomTask()
	ctx.JSON(200, gin.H{
		"task_id": task,
	})
}

func generateRandomTask() string {
	return "task" + fmt.Sprint(rand.Intn(1000))
}
