package main

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

type postStruct struct {
	Classname    string
	Teacher_name string
	Project_name string
	Class_time   string
	Status       string
	Stu_num      int
	Student_name int
	Use_time     string
}

type cellValue struct {
	sheet string
	cell  string
	value string
}
