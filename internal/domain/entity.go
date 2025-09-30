package domain

import (
	"github.com/gofrs/uuid/v5"
	"time"
)

type Item struct {
	ID        uuid.UUID
	Sort      int64
	Name      string
	Related   []Nested
	CreatedAt time.Time
	UpdatedAt *time.Time
}

func (it Item) Empty() bool {
	return it.ID == uuid.Nil
}

type Nested struct {
	ID      uuid.UUID
	Name    string
	Sort    int64
	Related []Atom
}

type Atom struct {
	ID   uuid.UUID
	Name string
}

// pagination & sorting

type Pagination struct {
	Limit  int
	Offset int
}

type SortOrder string

const (
	_         SortOrder = "ASC"
	OrderDesc SortOrder = "DESC"
)
