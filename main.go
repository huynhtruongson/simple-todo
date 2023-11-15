package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/huynhtruongson/simple-todo/middleware"
	auth_port "github.com/huynhtruongson/simple-todo/services/auth/port/http"
	task_port "github.com/huynhtruongson/simple-todo/services/task/port/http"
	user_port "github.com/huynhtruongson/simple-todo/services/user/port/http"
	"github.com/huynhtruongson/simple-todo/token"
	"github.com/huynhtruongson/simple-todo/utils"

	"github.com/jackc/pgx/v5"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load configuration", err)
	}
	db, err := pgx.Connect(context.Background(), config.DBAddress)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	}
	defer db.Close(context.Background())
	tokenMaker, err := token.NewPasetoMaker(config.TokenKey)
	if err != nil {
		log.Fatal("cannot initialize token maker", err)
	}
	r := gin.Default()
	v1 := r.Group("/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/login", auth_port.Login(db, tokenMaker))
			auth.POST("/renew-token", auth_port.RenewToken(db, tokenMaker))
		}

		user := v1.Group("/user")
		user.Use(middleware.AuthMiddleware(tokenMaker))
		{
			user.POST("/create", user_port.CreateUser(db))
		}

		task := v1.Group("/task")
		task.Use(middleware.AuthMiddleware(tokenMaker))
		{
			task.GET("/list", task_port.ListTask(db))
			task.POST("/create", task_port.CreateTask(db))
			task.PUT("/update/:id", task_port.UpdateTask(db))
			task.PUT("/delete/:id", task_port.DeleteTask(db))
		}
	}
	r.Run(config.ServerPort)
}
