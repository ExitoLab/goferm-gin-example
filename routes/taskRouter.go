package routes

import (
	controller "goferm-gin-example/controllers"

	"github.com/gin-gonic/gin"
)

//TaskRouter function
func TaskRouter(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/api/v1/tasks", controller.AddTask())
	incomingRoutes.GET("/api/v1/tasks/:task_id", controller.GetTaskByTaskID())
	incomingRoutes.GET("/api/v1/tasks", controller.GetAllTasks())
	incomingRoutes.DELETE("/api/v1/tasks/:task_id", controller.DeleteTaskByTaskID())
	incomingRoutes.PUT("/api/v1/tasks/:task_id", controller.UpdateTaskByTaskID())
	incomingRoutes.GET("/health", controller.HealthCheck())
}
