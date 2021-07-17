package models

import (
	"time"

	"gorm.io/gorm"
)

type Record struct {
	gorm.Model
	TeacherName   string `gorm:"not null"`
	ClassName     string
	ProjectName   string
	StudentNumber int
	StudentNum    int
	UseDate       time.Time
	ClassTime     string
	Status        string
}

type PostStruct struct {
	Classname   string `json:"classname"`
	TeacherName string `json:"teacher_name"`
	ProjectName string `json:"project_name"`
	ClassTime   string `json:"class_time"`
	Status      string `json:"status"`
	StuNum      int    `json:"stu_num" binding:"numeric"`
	StudentNum  int    `json:"student_num" binding:"numeric"`
	UseTime     string `json:"use_time"`
}

type CellValue struct {
	Sheet string
	Cell  string
	Value string
}
