// Package pageable provides offset-based and cursor-based pagination for Go REST APIs.
package pageable

const (
	// DefaultPage is the default page number (1-indexed).
	DefaultPage = 1
	// DefaultSize is the default number of items per page.
	DefaultSize = 10
	// MaxSize is the maximum allowed page size.
	MaxSize = 1000

	// DefaultCursorSize is the default number of items for cursor-based pagination.
	DefaultCursorSize = 10
	// MaxCursorSize is the maximum allowed cursor page size.
	MaxCursorSize = 1000
)
