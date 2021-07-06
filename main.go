package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	excelize "github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gin-gonic/gin"
	"github.com/xlzd/gotp"
)

const SerectCode = "DKEIR5BYLXTECP7BLI2C4WIUPKGFOAGE"

func init() {
	// 第一次创建数据库
	dir, _ := os.Getwd()
	databaseFile := filepath.Join(dir, "database.sqlite")
	if ok := existPath(databaseFile); !ok {
		db := getConn()
		_ = db.AutoMigrate(&Record{})
	}
	LogPath := filepath.Join(dir, "logs")
	if ok := existPath(LogPath); !ok {
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

	r.POST("/data", insertData)
	r.GET("/download", getDataFromDate)
	_ = r.Run(":8090")
}

// 插入数据库
func insertData(c *gin.Context) {
	// stc := make(map[string]interface{})
	// c.BindJSON(&stc)
	json := postStruct{}
	err := c.BindJSON(&json)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "请确认传值内容"})
		return
	}
	className := json.Classname
	teacherName := json.Teacher_name
	projectName := json.Project_name
	classTime := json.Class_time
	status := json.Status
	studentNumber := json.Stu_num
	studentNum := json.Student_name
	useTime := json.Use_time

	if className == "" || teacherName == "" || useTime == "" || classTime == "" {
		msg := "classname、teacher_name、use_time、class_time必须传值"
		c.JSON(http.StatusBadRequest, gin.H{"msg": msg})
		return
	} else {
		useTime, err := time.Parse("2006-01-02", useTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "use_time参数错误"})
			return
		}
		record := Record{
			TeacherName:   teacherName,
			ClassName:     className,
			ProjectName:   projectName,
			StudentNumber: studentNumber,
			StudentNum:    studentNum,
			UseDate:       useTime,
			ClassTime:     classTime,
			Status:        status,
		}
		db := getConn()
		db.Create(&record)
		c.JSON(http.StatusOK, gin.H{"msg": "插入成功！"})
	}
}

// 读取相关数据
func getDataFromDate(c *gin.Context) {
	timeFormat := "2006-01-02"
	timeUnix := int(time.Now().Unix())
	startDate := c.DefaultQuery("start_time", "2000-01-01")
	endDate := c.DefaultQuery("end_time", time.Now().Format(timeFormat))
	screct := c.DefaultQuery("screct", "")
	totp := gotp.NewDefaultTOTP(SerectCode)
	ok := totp.Verify(screct, timeUnix)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "验证码错误"})
		return
	}
	startDate1, err1 := time.Parse(timeFormat, startDate)
	endDate1, err2 := time.Parse(timeFormat, endDate)
	if err1 != nil || err2 != nil || !startDate1.Before(endDate1) {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "日期格式错误或时间前后错误，请检查日期格式"})
		return
	}
	db := getConn()
	records := []Record{}
	result := db.Where("use_date BETWEEN ? AND ?", startDate, endDate).Find(&records)
	if result.RowsAffected != 0 {
		f := excelize.NewFile()
		index := f.NewSheet("数据统计")
		f.SetActiveSheet(index)
		cellValues := make([]*cellValue, 0)
		cellValues = append(cellValues, &cellValue{
			sheet: "数据统计",
			cell:  "A1",
			value: "上课时间",
		}, &cellValue{
			sheet: "数据统计",
			cell:  "B1",
			value: "课程节次",
		}, &cellValue{
			sheet: "数据统计",
			cell:  "C1",
			value: "任课老师",
		}, &cellValue{
			sheet: "数据统计",
			cell:  "D1",
			value: "班级名称",
		}, &cellValue{
			sheet: "数据统计",
			cell:  "E1",
			value: "项目名称",
		}, &cellValue{
			sheet: "数据统计",
			cell:  "F1",
			value: "计划人数",
		}, &cellValue{
			sheet: "数据统计",
			cell:  "G1",
			value: "实际人数",
		}, &cellValue{
			sheet: "数据统计",
			cell:  "H1",
			value: "设备状态",
		})

		rowNum := 1
		// 创建标题栏
		for _, cellValue := range cellValues {
			err := f.SetCellValue(cellValue.sheet, cellValue.cell, cellValue.value)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "生成Excel出错。",
				})
				return
			}
		}
		// 插入数据
		for _, record := range records {
			rowNum++
			for k, v := range cellValues {
				var A rune = 'A'
				v.cell = fmt.Sprintf("%v%v", string(A+rune(k)), rowNum)
				switch k {
				case 0:
					v.value = record.UseDate.Format(timeFormat)
				case 1:
					v.value = record.ClassTime
				case 2:
					v.value = record.TeacherName
				case 3:
					v.value = record.ClassName
				case 4:
					v.value = record.ProjectName
				case 5:
					v.value = strconv.Itoa(record.StudentNum)
				case 6:
					v.value = strconv.Itoa(record.StudentNum)
				case 7:
					v.value = record.Status
				}
			}
			for _, data := range cellValues {
				err := f.SetCellValue(data.sheet, data.cell, data.value)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"msg": "生成Excel出错。",
					})
					return
				}
			}
		}
		c.Header("Content-Disposition", "attachment; filename=结果文件.xlsx")
		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		f.Write(c.Writer)
	} else if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": result.Error.Error()})

	} else {
		c.JSON(http.StatusNotFound, gin.H{"msg": "没有找到任何数据。"})

	}

}
