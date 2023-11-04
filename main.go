package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	task_port "github.com/sondev/todo-list/services/task/port/http"
	user_port "github.com/sondev/todo-list/services/user/port/http"
)

func main() {
	r := gin.Default()
	db, err := pgx.Connect(context.Background(), "postgresql://postgres:admin@localhost:5432/simple_todo?sslmode=disable")

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	}
	defer db.Close(context.Background())

	v1 := r.Group("/v1")
	{
		user := v1.Group("/user")
		{
			user.POST("/create", user_port.CreateUser(db))
		}
		task := v1.Group("/task")
		{
			task.GET("/list", task_port.ListTask(db))
			task.POST("/create", task_port.CreateTask(db))
			task.PUT("/update/:id", task_port.UpdateTask(db))
			task.PUT("/delete/:id", task_port.DeleteTask(db))
		}
	}
	r.Run(":3000")
}
