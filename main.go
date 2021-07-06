package main

import (
	"io"
	"log"
	"os"
	"path/filepath"

	handle "backend/src/handles"
	structs "backend/src/models"
	util "backend/src/utils"

	"github.com/gin-gonic/gin"
)

func init() {
	// 第一次创建数据库
	dir, _ := os.Getwd()
	databaseFile := filepath.Join(dir, "database.sqlite")
	if ok := util.ExistPath(databaseFile); !ok {
		db := util.GetConn()
		_ = db.AutoMigrate(&structs.Record{})
	}
	LogPath := filepath.Join(dir, "logs")
	if ok := util.ExistPath(LogPath); !ok {
		_ = os.Mkdir(LogPath, 0777)
	}
}

func main() {
	gin.SetMode(gin.ReleaseMode)

	logfile, err := os.OpenFile("./Logs/runtime.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to create request log file:", err)
	}

	errlogfile, err := os.OpenFile("./Logs/error.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to create request log file:", err)
	}

	gin.DefaultWriter = io.MultiWriter(logfile)
	gin.DefaultErrorWriter = io.MultiWriter(errlogfile)

	r := gin.Default()

	r.POST("/data", handle.InsertData)
	r.GET("/download", handle.GetDataFromDate)
	_ = r.Run(":8090")
}
