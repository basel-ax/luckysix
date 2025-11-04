package v1

import (
	"net/http"

	"github.com/basel-ax/luckysix/internal/service"
	"github.com/basel-ax/luckysix/pkg/logger"
	"github.com/gin-gonic/gin"
)

type tasksRoutes struct {
	t service.Tasks
	l *logger.Logger
}

func newTasksRoutes(handler *gin.RouterGroup, t service.Tasks, l *logger.Logger) {
	r := &tasksRoutes{t, l}

	handler.POST("/tasks/run", r.runTask)
}

func (r *tasksRoutes) runTask(c *gin.Context) {
	taskName := c.Query("name")
	if taskName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "task name is required"})
		return
	}

	err := r.t.RunTask(taskName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "task started"})
}
