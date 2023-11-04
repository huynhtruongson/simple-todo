package common

import "time"

type SQLModel struct {
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	DeletedAt time.Time `json:"-"`
}

type Paging struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}
