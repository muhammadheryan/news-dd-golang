package domain

import "time"

type News struct {
	Id         int
	Title      string `json:"title"`
	Content    string `json:"content"`
	Status     string `json:"status"`
	Topic      int    `json:"topic"`
	Tags       string `json:"tags"`
	CreateTime time.Time
	UpdateTime time.Time
}

type NewsRespond struct {
	Id         int
	Title      string
	Content    string
	Status     string
	Topic      string
	Tags       string
	CreateTime time.Time
	UpdateTime time.Time
}
