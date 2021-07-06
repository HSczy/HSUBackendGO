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
	Classname   string
	TeacherName string
	ProjectName string
	ClassTime   string
	Status      string
	StuNum      int
	StudentName int
	UseTime     string
}

type CellValue struct {
	Sheet string
	Cell  string
	Value string
}
