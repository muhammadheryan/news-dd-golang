package domain

import "time"

type Topic struct {
	Id         int    `json:"id"`
	Topic      string `json:"topic"`
	CreateTime time.Time
	UpdateTime time.Time
}
