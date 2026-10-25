package model

import "uuid"

type Car struct {
	ID    uuid.UUID
	Brand string
	Model string
	Year  int
}

type Pagination struct {
	PageSize  int32
	PageToken uuid.UUID
}

type CarFilter struct {
	Search *string `json:"search"`
	Year   *int    `json:"year"`
}

type CarColumnSort string
type SortDirection string

const (
	SortDirectionAsc  SortDirection = "asc"
	SortDirectionDesc SortDirection = "desc"
)

const (
	CarColumnSortID    CarColumnSort = "id"
	CarColumnSortBrand CarColumnSort = "brand"
	CarColumnSortModel CarColumnSort = "model"
	CarColumnSortYear  CarColumnSort = "year"
)

type CarSortBy struct {
	Field     CarColumnSort `json:"field"`
	Direction SortDirection `json:"direction"`
}
