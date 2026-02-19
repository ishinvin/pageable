package pageable

// CursorPageMetadata holds pagination metadata for cursor-based pagination.
// Unlike PageMetadata, it does not include TotalItems or TotalPages, since
// cursor-based pagination avoids COUNT queries for better performance.
type CursorPageMetadata struct {
	NextCursor string `json:"nextCursor"`
	PrevCursor string `json:"prevCursor"`
	HasNext    bool   `json:"hasNext"`
	HasPrev    bool   `json:"hasPrev"`
	Size       int    `json:"size"`
}

// CursorPage represents a paginated response for cursor-based pagination.
type CursorPage[T any] struct {
	Items    []T                `json:"items"`
	Metadata CursorPageMetadata `json:"metadata"`
}

// EmptyCursorPage creates an empty CursorPage with no items and no cursors.
func EmptyCursorPage[T any](size int) CursorPage[T] {
	return NewCursorPage[T](nil, "", "", false, false, size)
}

// NewCursorPage creates a CursorPage from items and cursor information.
// The nextCursor and prevCursor should be pre-encoded cursor strings
// (use EncodeCursor to produce them).
// A nil items slice is converted to an empty slice to ensure JSON serializes as [] not null.
func NewCursorPage[T any](items []T, nextCursor, prevCursor string, hasNext, hasPrev bool, size int) CursorPage[T] {
	if items == nil {
		items = make([]T, 0)
	}
	return CursorPage[T]{
		Items: items,
		Metadata: CursorPageMetadata{
			NextCursor: nextCursor,
			PrevCursor: prevCursor,
			HasNext:    hasNext,
			HasPrev:    hasPrev,
			Size:       size,
		},
	}
}
