package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/xlzd/gotp"

	structs "backend/src/models"
	util "backend/src/utils"

	"github.com/gin-gonic/gin"
)

const SecretCode = "DKEIR5BYLXTECP7BLI2C4WIUPKGFOAGE"

// InsertData 插入数据库
func InsertData(c *gin.Context) {
	// stc := make(map[string]interface{})
	// c.BindJSON(&stc)
	json := structs.PostStruct{}
	err := c.BindJSON(&json)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "请确认传值内容"})
		return
	}
	className := json.Classname
	teacherName := json.TeacherName
	projectName := json.ProjectName
	classTime := json.ClassTime
	status := json.Status
	studentNumber := json.StuNum
	studentNum := json.StudentNum
	useTime := json.UseTime

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
		record := structs.Record{
			TeacherName:   teacherName,
			ClassName:     className,
			ProjectName:   projectName,
			StudentNumber: studentNumber,
			StudentNum:    studentNum,
			UseDate:       useTime,
			ClassTime:     classTime,
			Status:        status,
		}
		db := util.GetConn()
		db.Create(&record)
		c.JSON(http.StatusOK, gin.H{"msg": "插入成功！"})
	}
}

// GetDataFromDate 读取相关数据
func GetDataFromDate(c *gin.Context) {
	timeFormat := "2006-01-02"
	timeUnix := int(time.Now().Unix())
	startDate := c.DefaultQuery("start_time", "2000-01-01")
	endDate := c.DefaultQuery("end_time", time.Now().Format(timeFormat))
	secret := c.DefaultQuery("secret", "")
	totp := gotp.NewDefaultTOTP(SecretCode)
	ok := totp.Verify(secret, timeUnix)
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
	db := util.GetConn()
	var records []structs.Record
	result := db.Where("use_date BETWEEN ? AND ?", startDate, endDate).Find(&records)
	if result.RowsAffected != 0 {
		f := excelize.NewFile()
		index := f.NewSheet("数据统计")
		f.SetActiveSheet(index)
		cellValues := make([]*structs.CellValue, 0)
		cellValues = append(cellValues, &structs.CellValue{
			Sheet: "数据统计",
			Cell:  "A1",
			Value: "上课时间",
		}, &structs.CellValue{
			Sheet: "数据统计",
			Cell:  "B1",
			Value: "课程节次",
		}, &structs.CellValue{
			Sheet: "数据统计",
			Cell:  "C1",
			Value: "任课老师",
		}, &structs.CellValue{
			Sheet: "数据统计",
			Cell:  "D1",
			Value: "班级名称",
		}, &structs.CellValue{
			Sheet: "数据统计",
			Cell:  "E1",
			Value: "项目名称",
		}, &structs.CellValue{
			Sheet: "数据统计",
			Cell:  "F1",
			Value: "计划人数",
		}, &structs.CellValue{
			Sheet: "数据统计",
			Cell:  "G1",
			Value: "实际人数",
		}, &structs.CellValue{
			Sheet: "数据统计",
			Cell:  "H1",
			Value: "设备状态",
		})

		rowNum := 1
		// 创建标题栏
		for _, cellValue := range cellValues {
			err := f.SetCellValue(cellValue.Sheet, cellValue.Cell, cellValue.Value)
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
				v.Cell = fmt.Sprintf("%v%v", string(A+rune(k)), rowNum)
				switch k {
				case 0:
					v.Value = record.UseDate.Format(timeFormat)
				case 1:
					v.Value = record.ClassTime
				case 2:
					v.Value = record.TeacherName
				case 3:
					v.Value = record.ClassName
				case 4:
					v.Value = record.ProjectName
				case 5:
					v.Value = strconv.Itoa(record.StudentNum)
				case 6:
					v.Value = strconv.Itoa(record.StudentNum)
				case 7:
					v.Value = record.Status
				}
			}
			for _, data := range cellValues {
				err := f.SetCellValue(data.Sheet, data.Cell, data.Value)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"msg": "生成Excel出错。",
					})
					return
				}
			}
		}
		c.Header("Content-Disposition", "attachment; filename=结果文件.xlsx")
		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.Sheet")
		err := f.Write(c.Writer)
		if err != nil {
			return 
		}
	} else if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": result.Error.Error()})

	} else {
		c.JSON(http.StatusNotFound, gin.H{"msg": "没有找到任何数据。"})

	}

}
