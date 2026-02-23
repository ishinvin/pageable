package pageable

import (
	"net/url"
	"strconv"
	"strings"
)

// CursorRequest represents cursor-based pagination parameters.
type CursorRequest struct {
	Cursor string
	Size   int
	Sort   []Sort
}

// NewCursorRequest creates a CursorRequest with defaults applied.
// Size is clamped to [1, MaxCursorSize].
func NewCursorRequest(cursor string, size int, sort []Sort) CursorRequest {
	if size < 1 {
		size = DefaultCursorSize
	}
	if size > MaxCursorSize {
		size = MaxCursorSize
	}
	return CursorRequest{Cursor: cursor, Size: size, Sort: sort}
}

// CursorRequestFromQuery parses a CursorRequest from URL query parameters.
// Recognized keys: "cursor", "size", "sort".
// Defaults: empty cursor (first page), DefaultCursorSize.
func CursorRequestFromQuery(values url.Values) CursorRequest {
	cursor := values.Get("cursor")

	size := DefaultCursorSize
	if v := values.Get("size"); v != "" {
		if s, err := strconv.Atoi(v); err == nil && s > 0 {
			size = s
		}
	}
	if size > MaxCursorSize {
		size = MaxCursorSize
	}

	var sort []Sort
	if sortParams := values["sort"]; len(sortParams) > 0 {
		sort = ParseSorts(sortParams)
	}

	return CursorRequest{Cursor: cursor, Size: size, Sort: sort}
}

// SortableFields filters sorts to only include the specified fields.
// Any sort with a field not in the allowed list is removed.
func (cr CursorRequest) SortableFields(fields ...string) CursorRequest {
	cr.Sort = filterSortsByFields(cr.Sort, fields...)
	return cr
}

// MapSortFields replaces sort field names using the provided mapping.
// Use this to translate user-facing field names (e.g., "createdAt") to
// database column names (e.g., "created_at"). Unmapped fields are kept as-is.
func (cr CursorRequest) MapSortFields(fieldMap map[string]string) CursorRequest {
	cr.Sort = mapSortFields(cr.Sort, fieldMap)
	return cr
}

// WithDefaultSort sets the sort to the given defaults if no sort is set.
// Has no effect if the request already has sorts from query parameters.
func (cr CursorRequest) WithDefaultSort(sorts ...Sort) CursorRequest {
	if cr.Sort == nil {
		cr.Sort = sorts
	}
	return cr
}

// OrderBy returns an ORDER BY clause string from the request's sorts.
// Returns a string like "name desc, id asc".
// Returns an empty string if no sorts are set.
func (cr CursorRequest) OrderBy() string {
	if len(cr.Sort) == 0 {
		return ""
	}
	parts := make([]string, len(cr.Sort))
	for i, s := range cr.Sort {
		parts[i] = s.Field + " " + string(s.Direction)
	}
	return strings.Join(parts, ", ")
}

// Limit returns Size + 1 for database queries.
// Querying one extra item is the standard way to detect whether more items exist
// without a separate COUNT query.
func (cr CursorRequest) Limit() int {
	return cr.Size + 1
}

// HasCursor returns true if a non-empty cursor was provided.
func (cr CursorRequest) HasCursor() bool {
	return cr.Cursor != ""
}

// DecodedCursor decodes and returns the full CursorData.
// Returns (CursorData{}, nil) if no cursor is set.
func (cr CursorRequest) DecodedCursor() (CursorData, error) {
	if cr.Cursor == "" {
		return CursorData{}, nil
	}
	return DecodeCursor(cr.Cursor)
}
