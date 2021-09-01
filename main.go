package main

import (
	"backend/src/middleware"
	"io"
	"log"
	"os"
	"path/filepath"

	"backend/src/controllers"
	"backend/src/models"
	"backend/src/utils"

	"github.com/gin-gonic/gin"
)

func init() {
	// 第一次创建数据库
	dir, _ := os.Getwd()
	databaseFile := filepath.Join(dir, "database.sqlite")
	if ok := utils.ExistPath(databaseFile); !ok {
		db := utils.GetConn()
		_ = db.AutoMigrate(&models.Record{})
	}
	LogPath := filepath.Join(dir, "Logs")
	if ok := utils.ExistPath(LogPath); !ok {
		_ = os.Mkdir(LogPath, 0777)
	}
}

func main() {
	//gin.SetMode(gin.ReleaseMode)

	logfile, err := os.OpenFile("./Logs/runtime.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to create request log file:", err)
	}

	errlogfile, err := os.OpenFile("./Logs/error.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to create request log file:", err)
	}

	gin.DefaultWriter = io.MultiWriter(logfile, os.Stdout)
	gin.DefaultErrorWriter = io.MultiWriter(errlogfile)

	r := gin.Default()
	r.Use(middleware.Cors())
	r.POST("/data", controller.InsertData)
	r.GET("/download", controller.GetDataFromDate)
	_ = r.Run(":8090")
}
