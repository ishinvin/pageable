package pageable

import (
	"net/url"
	"strconv"
	"strings"
)

// PageRequest represents offset-based pagination parameters.
type PageRequest struct {
	Page int
	Size int
	Sort []Sort
}

// NewPageRequest creates a PageRequest with defaults applied.
// Page is clamped to a minimum of 1. Size is clamped to [1, MaxSize].
func NewPageRequest(page, size int, sort []Sort) PageRequest {
	if page < 1 {
		page = DefaultPage
	}
	if size < 1 {
		size = DefaultSize
	}
	if size > MaxSize {
		size = MaxSize
	}
	return PageRequest{Page: page, Size: size, Sort: sort}
}

// PageRequestFromQuery parses a PageRequest from URL query parameters.
// Recognized keys: "page", "size", "sort".
// Uses DefaultPage and DefaultSize for missing or invalid values.
// Size is clamped to [1, MaxSize].
func PageRequestFromQuery(values url.Values) PageRequest {
	page := DefaultPage
	if v := values.Get("page"); v != "" {
		if p, err := strconv.Atoi(v); err == nil && p > 0 {
			page = p
		}
	}

	size := DefaultSize
	if v := values.Get("size"); v != "" {
		if s, err := strconv.Atoi(v); err == nil && s > 0 {
			size = s
		}
	}
	if size > MaxSize {
		size = MaxSize
	}

	var sort []Sort
	if sortParams := values["sort"]; len(sortParams) > 0 {
		sort = ParseSorts(sortParams)
	}

	return PageRequest{Page: page, Size: size, Sort: sort}
}

// Offset returns the zero-based offset for database queries.
// Calculated as (Page - 1) * Size.
func (pr PageRequest) Offset() int {
	return (pr.Page - 1) * pr.Size
}

// Limit returns the page size (alias for clarity in SQL-style usage).
func (pr PageRequest) Limit() int {
	return pr.Size
}

// SortableFields filters sorts to only include the specified fields.
// Any sort with a field not in the allowed list is removed.
func (pr PageRequest) SortableFields(fields ...string) PageRequest {
	pr.Sort = filterSortsByFields(pr.Sort, fields...)
	return pr
}

// WithDefaultSort sets the sort to the given defaults if no sort is set.
// Has no effect if the request already has sorts from query parameters.
func (pr PageRequest) WithDefaultSort(sorts ...Sort) PageRequest {
	if pr.Sort == nil {
		pr.Sort = sorts
	}
	return pr
}

// OrderBy returns an ORDER BY clause string from the request's sorts.
// Returns a string like "name desc, id asc".
// Returns an empty string if no sorts are set.
func (pr PageRequest) OrderBy() string {
	if len(pr.Sort) == 0 {
		return ""
	}
	parts := make([]string, len(pr.Sort))
	for i, s := range pr.Sort {
		parts[i] = s.Field + " " + string(s.Direction)
	}
	return strings.Join(parts, ", ")
}
