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

type Nested struct {
	ID      uuid.UUID
	Name    string
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
	OrderAsc  SortOrder = "ASC"
	OrderDesc SortOrder = "DESC"
)
