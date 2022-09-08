package models

type Pagination struct {
	page     int `json:"page"`
	pageSize int `json:"pageSize"`
}
