package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	task_port "github.com/huynhtruongson/simple-todo/services/task/port/http"
	user_port "github.com/huynhtruongson/simple-todo/services/user/port/http"
	"github.com/huynhtruongson/simple-todo/utils"

	"github.com/jackc/pgx/v5"
)

func main() {
	r := gin.Default()
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load configuration", err)
	}
	db, err := pgx.Connect(context.Background(), config.DBAddress)

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
	r.Run(config.ServerPort)
}
