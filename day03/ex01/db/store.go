package db

import (
	"day03/ex01/types"
)

type Store interface {
	GetPlaces(limit int, offset int) ([]types.Place, int, error)
}
