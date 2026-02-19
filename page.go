package pageable

// PageMetadata holds pagination metadata for offset-based pagination.
type PageMetadata struct {
	Page       int   `json:"page"`
	Size       int   `json:"size"`
	TotalItems int64 `json:"totalItems"`
	TotalPages int   `json:"totalPages"`
}

// Page represents a paginated response for offset-based pagination.
type Page[T any] struct {
	Items    []T          `json:"items"`
	Metadata PageMetadata `json:"metadata"`
}

// EmptyPage creates an empty Page with zero results, preserving the request's page and size.
func EmptyPage[T any](request PageRequest) Page[T] {
	return NewPage[T](nil, request, 0)
}

// NewPage creates a Page from items, request parameters, and total item count.
// TotalPages is calculated automatically via ceiling division.
// A nil items slice is converted to an empty slice to ensure JSON serializes as [] not null.
func NewPage[T any](items []T, request PageRequest, totalItems int64) Page[T] {
	if items == nil {
		items = make([]T, 0)
	}

	totalPages := 0
	if totalItems > 0 && request.Size > 0 {
		totalPages = int((totalItems + int64(request.Size) - 1) / int64(request.Size))
	}

	return Page[T]{
		Items: items,
		Metadata: PageMetadata{
			Page:       request.Page,
			Size:       request.Size,
			TotalItems: totalItems,
			TotalPages: totalPages,
		},
	}
}
