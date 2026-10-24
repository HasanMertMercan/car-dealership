package model

import "uuid"

type Car struct {
	ID    uuid.UUID
	Brand string
	Model string
	Year  int
}

type CarFilter struct {
	Search *string `json:"search"`
	Year   *int    `json:"year"`
}
