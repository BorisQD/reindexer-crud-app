package domain

import "strings"

// Transform реализует условную специфичную для домена логику обработки сущности с исключением части вложенных документов.
// Предполагается что она не возвращает ошибок
func Transform(item Item) Item {
	var filtered []Nested
	for _, nst := range item.Related {
		if strings.HasPrefix(nst.Name, "deprecated_") {
			continue
		}

		var filteredAtoms []Atom
		for _, atom := range nst.Related {
			if strings.HasPrefix(atom.Name, "deleted_") {
				continue
			}

			filteredAtoms = append(filteredAtoms, atom)
		}
		nst.Related = filteredAtoms

		filtered = append(filtered, nst)
	}
	item.Related = filtered

	return item
}
