package pageable

import (
	"strings"
)

// Direction represents the sort direction.
type Direction string

const (
	// ASC sorts in ascending order.
	ASC Direction = "asc"
	// DESC sorts in descending order.
	DESC Direction = "desc"
)

// Sort represents a single sort field with its direction.
type Sort struct {
	Field     string
	Direction Direction
}

// String returns the sort as "field,direction" (e.g., "name,desc" or "id,asc").
func (s Sort) String() string {
	return s.Field + "," + string(s.Direction)
}

// ParseSort parses a "field,direction" sort string into a Sort.
// Format: "field,direction" where direction is "asc" or "desc".
// If direction is omitted, defaults to ascending.
// Examples: "id,desc" -> {id, desc}, "name,asc" -> {name, asc}, "name" -> {name, asc}.
// Returns nil for empty input.
func ParseSort(raw string) *Sort {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}

	parts := strings.SplitN(raw, ",", 2)
	field := strings.TrimSpace(parts[0])
	if field == "" || !isSafeIdentifier(field) {
		return nil
	}

	dir := ASC
	if len(parts) == 2 {
		d := strings.TrimSpace(strings.ToLower(parts[1]))
		if Direction(d) == DESC {
			dir = DESC
		}
	}

	return &Sort{Field: field, Direction: dir}
}

// ParseSorts parses multiple "field,direction" sort strings as typically received
// from url.Values where ?sort=id,desc&sort=name,asc yields []string{"id,desc", "name,asc"}.
// Each string is parsed with ParseSort.
func ParseSorts(raw []string) []Sort {
	var sorts []Sort
	for _, r := range raw {
		if s := ParseSort(r); s != nil {
			sorts = append(sorts, *s)
		}
	}
	if len(sorts) == 0 {
		return nil
	}
	return sorts
}

// filterSortsByFields returns only sorts whose field is in the allowed list.
// Returns nil if no sorts match.
func filterSortsByFields(sorts []Sort, fields ...string) []Sort {
	allowed := make(map[string]struct{}, len(fields))
	for _, f := range fields {
		allowed[f] = struct{}{}
	}
	var filtered []Sort
	for _, s := range sorts {
		if _, ok := allowed[s.Field]; ok {
			filtered = append(filtered, s)
		}
	}
	if len(filtered) == 0 {
		return nil
	}
	return filtered
}

// isSafeIdentifier checks that a field name contains only safe SQL identifier
// characters: letters, digits, underscores, and dots (for table-qualified names like "posts.id").
func isSafeIdentifier(s string) bool {
	for _, c := range s {
		if (c < 'a' || c > 'z') && (c < 'A' || c > 'Z') && (c < '0' || c > '9') && c != '_' && c != '.' {
			return false
		}
	}
	return true
}
