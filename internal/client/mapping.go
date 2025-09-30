package client

import (
	"crud/internal/domain"
	"github.com/gofrs/uuid/v5"
	"time"
)

type Item struct {
	ID        string     `reindex:"id,,pk"`
	Sort      int64      `reindex:"sort"`
	Name      string     `reindex:"name"`
	Related   []Nested   `reindex:"related"`
	CreatedAt time.Time  `reindex:"createdAt"`
	UpdatedAt *time.Time `reindex:"updatedAt"`
}

func (it Item) toModel() domain.Item {
	related := make([]domain.Nested, 0, len(it.Related))
	for _, a := range it.Related {
		related = append(related, a.toModel())
	}

	id, _ := uuid.FromString(it.ID)

	return domain.Item{
		ID:        id,
		Sort:      it.Sort,
		Name:      it.Name,
		Related:   related,
		CreatedAt: it.CreatedAt,
		UpdatedAt: it.UpdatedAt,
	}
}

type Nested struct {
	ID      string `reindex:"id"`
	Name    string `reindex:"name"`
	Sort    int64  `reindex:"sort"`
	Related []Atom `reindex:"related"`
}

func (nst Nested) toModel() domain.Nested {
	atoms := make([]domain.Atom, 0, len(nst.Related))
	for _, a := range nst.Related {
		atoms = append(atoms, a.toModel())
	}

	id, _ := uuid.FromString(nst.ID)

	return domain.Nested{
		ID:      id,
		Name:    nst.Name,
		Sort:    nst.Sort,
		Related: atoms,
	}
}

func nestedToDTO(n domain.Nested) Nested {
	return Nested{
		ID:      n.ID.String(),
		Name:    n.Name,
		Sort:    n.Sort,
		Related: nil,
	}
}

type Atom struct {
	ID   string `reindex:"id"`
	Name string `reindex:"name"`
}

func (a Atom) toModel() domain.Atom {
	id, _ := uuid.FromString(a.ID)
	return domain.Atom{
		ID:   id,
		Name: a.Name,
	}
}

func toDTO(it domain.Item) Item {
	subItems := make([]Nested, 0, len(it.Related))
	for _, nst := range it.Related {
		subItems = append(subItems, nestedToDTO(nst))
	}

	return Item{
		ID:        it.ID.String(),
		Sort:      it.Sort,
		Name:      it.Name,
		Related:   subItems,
		CreatedAt: it.CreatedAt,
		UpdatedAt: it.UpdatedAt,
	}
}
