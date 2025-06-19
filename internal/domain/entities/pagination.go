package entities

type PaginationRes[T any] struct {
	Data  []T
	Total int    `json:"total"`
	Prev  string `json:"prev"`
	Next  string `json:"next"`
}
