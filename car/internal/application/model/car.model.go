package model

import "uuid"

type Car struct {
	ID    uuid.UUID
	Brand string
	Model string
	Year  int
}
