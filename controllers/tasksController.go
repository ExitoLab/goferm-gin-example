package controllers

import (
	"context"
	"log"
	"strconv"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"goferm-gin-example/database"

	"goferm-gin-example/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var taskCollection *mongo.Collection = database.OpenCollection(database.Client, "task")

var validate = validator.New()

//check health endpoint
func HealthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"data":         "",
			"responseCode": http.StatusOK,
			"message":      "The application is running sucessfully",
		})
	}
}

//Add Task
func AddTask() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var task models.Tasks

		defer cancel()
		if err := c.BindJSON(&task); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"data":         "",
				"responseCode": http.StatusBadRequest,
				"message":      err.Error(),
			})
			return
		}

		validationErr := validate.Struct(task)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"data":         "",
				"responseCode": http.StatusBadRequest,
				"message":      validationErr.Error(),
			})
			return
		}

		task.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		task.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		task.ID = primitive.NewObjectID()
		task.Task_id = task.ID.Hex()

		resultInsertionNumber, insertErr := taskCollection.InsertOne(ctx, task)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"data":         "",
				"responseCode": http.StatusInternalServerError,
				"message":      "task was not created",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data":         resultInsertionNumber,
			"responseCode": http.StatusOK,
			"message":      "The task has been created successfully",
		})

	}
}

//Get a single task by task_id
func GetTaskByTaskID() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		task_id := c.Param("task_id")

		var task models.Tasks

		err := taskCollection.FindOne(ctx, bson.M{"task_id": task_id}).Decode(&task)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"data":         "",
				"responseCode": http.StatusInternalServerError,
				"message":      "error occured while fetching task",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data":         task,
			"responseCode": http.StatusOK,
		})
	}
}

// Delete a Task
func DeleteTaskByTaskID() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		task_id := c.Param("task_id")
		defer cancel()

		result, err := taskCollection.DeleteOne(ctx, bson.M{"task_id": task_id})

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"data":         "",
				"responseCode": http.StatusBadRequest,
				"message":      "Error occured while trying to delete the document",
			})
			return
		}

		if result.DeletedCount > 0 {
			c.JSON(http.StatusOK, gin.H{
				"data":         "",
				"responseCode": http.StatusOK,
				"message":      "The task was successfully deleted",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"data":         "",
				"responseCode": http.StatusOK,
				"message":      "No Task to be deleted",
			})
		}
	}
}

//Update a Task
func UpdateTaskByTaskID() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var task models.Tasks

		task_id := c.Param("task_id")

		if err := c.BindJSON(&task); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"data":         "",
				"responseCode": http.StatusBadRequest,
				"message":      err.Error(),
			})
			return
		}

		var updateObj primitive.D

		if task.Title != nil {
			updateObj = append(updateObj, bson.E{"title", task.Title})
		}

		if task.Description != nil {
			updateObj = append(updateObj, bson.E{"description", task.Description})
		}

		task.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"created_at", task.Created_at})

		task.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"updated_at", task.Updated_at})

		upsert := true
		filter := bson.M{"task_id": task_id}
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := taskCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"$set", updateObj},
			},
			&opt,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"data":         "",
				"responseCode": http.StatusInternalServerError,
				"message":      err.Error(),
			})
			return
		}

		if result.ModifiedCount > 0 {
			c.JSON(http.StatusOK, gin.H{
				"data":         "",
				"responseCode": http.StatusOK,
				"message":      "The task has been successfully updated",
			})

		}
	}
}

//Get all tasks in the database
func GetAllTasks() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		// recordPerPage := 10
		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		page, err1 := strconv.Atoi(c.Query("page"))
		if err1 != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage
		startIndex, err = strconv.Atoi(c.Query("startIndex"))

		matchStage := bson.D{{"$match", bson.D{{}}}}
		groupStage := bson.D{{"$group", bson.D{{"_id", bson.D{{"_id", "null"}}}, {"total_count", bson.D{{"$sum", 1}}}, {"data", bson.D{{"$push", "$$ROOT"}}}}}}
		projectStage := bson.D{
			{"$project", bson.D{
				{"_id", 0},
				{"total_count", 1},
				{"total_tasks", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},
			}}}

		result, err := taskCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing tasks"})
		}
		var allTasks []bson.M
		if err = result.All(ctx, &allTasks); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allTasks[0])

	}
}
