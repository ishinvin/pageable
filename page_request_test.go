package pageable

import (
	"net/url"
	"testing"
)

func TestPageRequestFromQuery(t *testing.T) {
	tests := []struct {
		name         string
		values       url.Values
		expectedPage int
		expectedSize int
		expectedSort []Sort
	}{
		{
			name:         "defaults on empty values",
			values:       url.Values{},
			expectedPage: 1,
			expectedSize: 10,
			expectedSort: nil,
		},
		{
			name:         "valid page and size",
			values:       url.Values{"page": {"3"}, "size": {"25"}},
			expectedPage: 3,
			expectedSize: 25,
			expectedSort: nil,
		},
		{
			name:         "negative page clamps to default",
			values:       url.Values{"page": {"-1"}},
			expectedPage: 1,
			expectedSize: 10,
		},
		{
			name:         "zero page clamps to default",
			values:       url.Values{"page": {"0"}},
			expectedPage: 1,
			expectedSize: 10,
		},
		{
			name:         "size exceeds max",
			values:       url.Values{"size": {"5000"}},
			expectedPage: 1,
			expectedSize: MaxSize,
		},
		{
			name:         "non-numeric page",
			values:       url.Values{"page": {"abc"}},
			expectedPage: 1,
			expectedSize: 10,
		},
		{
			name:         "non-numeric size",
			values:       url.Values{"size": {"xyz"}},
			expectedPage: 1,
			expectedSize: 10,
		},
		{
			name:         "negative size clamps to default",
			values:       url.Values{"size": {"-5"}},
			expectedPage: 1,
			expectedSize: 10,
		},
		{
			name:         "with sort single param",
			values:       url.Values{"page": {"1"}, "sort": {"name,desc"}},
			expectedPage: 1,
			expectedSize: 10,
			expectedSort: []Sort{
				{Field: "name", Direction: DESC},
			},
		},
		{
			name:         "with sort multiple params",
			values:       url.Values{"sort": {"name,desc", "id,asc"}},
			expectedPage: 1,
			expectedSize: 10,
			expectedSort: []Sort{
				{Field: "name", Direction: DESC},
				{Field: "id", Direction: ASC},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := PageRequestFromQuery(tt.values)
			if req.Page != tt.expectedPage {
				t.Errorf("Page = %d, want %d", req.Page, tt.expectedPage)
			}
			if req.Size != tt.expectedSize {
				t.Errorf("Size = %d, want %d", req.Size, tt.expectedSize)
			}
			if tt.expectedSort == nil {
				if req.Sort != nil {
					t.Errorf("Sort = %v, want nil", req.Sort)
				}
			} else {
				if len(req.Sort) != len(tt.expectedSort) {
					t.Fatalf("Sort length = %d, want %d", len(req.Sort), len(tt.expectedSort))
				}
				for i, s := range req.Sort {
					if s.Field != tt.expectedSort[i].Field || s.Direction != tt.expectedSort[i].Direction {
						t.Errorf("Sort[%d] = %v, want %v", i, s, tt.expectedSort[i])
					}
				}
			}
		})
	}
}

func TestNewPageRequest(t *testing.T) {
	tests := []struct {
		name         string
		page, size   int
		expectedPage int
		expectedSize int
	}{
		{"valid", 2, 20, 2, 20},
		{"negative page", -1, 10, 1, 10},
		{"zero page", 0, 10, 1, 10},
		{"negative size", 1, -5, 1, 10},
		{"zero size", 1, 0, 1, 10},
		{"size over max", 1, 5000, 1, 1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := NewPageRequest(tt.page, tt.size, nil)
			if req.Page != tt.expectedPage {
				t.Errorf("Page = %d, want %d", req.Page, tt.expectedPage)
			}
			if req.Size != tt.expectedSize {
				t.Errorf("Size = %d, want %d", req.Size, tt.expectedSize)
			}
		})
	}
}

func TestPageRequestOffset(t *testing.T) {
	tests := []struct {
		page, size     int
		expectedOffset int
	}{
		{1, 10, 0},
		{2, 10, 10},
		{3, 20, 40},
		{1, 50, 0},
		{5, 25, 100},
	}

	for _, tt := range tests {
		req := PageRequest{Page: tt.page, Size: tt.size}
		if got := req.Offset(); got != tt.expectedOffset {
			t.Errorf("PageRequest{Page:%d, Size:%d}.Offset() = %d, want %d",
				tt.page, tt.size, got, tt.expectedOffset)
		}
	}
}

func TestPageRequestLimit(t *testing.T) {
	req := PageRequest{Page: 1, Size: 25}
	if got := req.Limit(); got != 25 {
		t.Errorf("Limit() = %d, want 25", got)
	}
}

func TestPageRequestSortableFields(t *testing.T) {
	tests := []struct {
		name     string
		sort     []Sort
		allowed  []string
		expected []Sort
	}{
		{
			name:     "keeps allowed fields",
			sort:     []Sort{{Field: "name", Direction: ASC}, {Field: "id", Direction: DESC}},
			allowed:  []string{"name", "id"},
			expected: []Sort{{Field: "name", Direction: ASC}, {Field: "id", Direction: DESC}},
		},
		{
			name:     "removes disallowed fields",
			sort:     []Sort{{Field: "name", Direction: ASC}, {Field: "secret", Direction: DESC}},
			allowed:  []string{"name", "id"},
			expected: []Sort{{Field: "name", Direction: ASC}},
		},
		{
			name:     "all fields disallowed",
			sort:     []Sort{{Field: "secret", Direction: ASC}},
			allowed:  []string{"name", "id"},
			expected: nil,
		},
		{
			name:     "no sorts",
			sort:     nil,
			allowed:  []string{"name"},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := PageRequest{Page: 1, Size: 10, Sort: tt.sort}
			req = req.SortableFields(tt.allowed...)
			if tt.expected == nil {
				if req.Sort != nil {
					t.Errorf("Sort = %v, want nil", req.Sort)
				}
			} else {
				if len(req.Sort) != len(tt.expected) {
					t.Fatalf("Sort length = %d, want %d", len(req.Sort), len(tt.expected))
				}
				for i, s := range req.Sort {
					if s.Field != tt.expected[i].Field || s.Direction != tt.expected[i].Direction {
						t.Errorf("Sort[%d] = %v, want %v", i, s, tt.expected[i])
					}
				}
			}
		})
	}
}

func TestPageRequestWithDefaultSort(t *testing.T) {
	// Applies default when no sort is set
	req := PageRequest{Page: 1, Size: 10}
	req = req.WithDefaultSort(Sort{Field: "id", Direction: ASC})
	if len(req.Sort) != 1 || req.Sort[0].Field != "id" {
		t.Errorf("Sort = %v, want [{id asc}]", req.Sort)
	}

	// Does not override existing sort
	req = PageRequest{Page: 1, Size: 10, Sort: []Sort{{Field: "name", Direction: DESC}}}
	req = req.WithDefaultSort(Sort{Field: "id", Direction: ASC})
	if len(req.Sort) != 1 || req.Sort[0].Field != "name" {
		t.Errorf("Sort = %v, want [{name desc}]", req.Sort)
	}
}

func TestPageRequestOrderBy(t *testing.T) {
	tests := []struct {
		name     string
		sort     []Sort
		expected string
	}{
		{"no sorts", nil, ""},
		{"single sort", []Sort{{Field: "name", Direction: ASC}}, "name asc"},
		{"multiple sorts", []Sort{
			{Field: "created_at", Direction: DESC},
			{Field: "id", Direction: ASC},
		}, "created_at desc, id asc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := PageRequest{Page: 1, Size: 10, Sort: tt.sort}
			if got := req.OrderBy(); got != tt.expected {
				t.Errorf("OrderBy() = %q, want %q", got, tt.expected)
			}
		})
	}
}
