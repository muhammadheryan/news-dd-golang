package domain

import "time"

type Tags struct {
	Id         int    `json:"id"`
	Tag        string `json:"tags"`
	CreateTime time.Time
	UpdateTime time.Time
}
