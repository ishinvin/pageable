package pageable

import (
	"net/url"
	"testing"
)

func TestCursorRequestFromQuery(t *testing.T) {
	tests := []struct {
		name           string
		values         url.Values
		expectedCursor string
		expectedSize   int
		expectedSort   []Sort
	}{
		{
			name:           "defaults on empty values",
			values:         url.Values{},
			expectedCursor: "",
			expectedSize:   DefaultCursorSize,
			expectedSort:   nil,
		},
		{
			name:           "with cursor and size",
			values:         url.Values{"cursor": {"abc123"}, "size": {"50"}},
			expectedCursor: "abc123",
			expectedSize:   50,
		},
		{
			name:           "size exceeds max",
			values:         url.Values{"size": {"5000"}},
			expectedCursor: "",
			expectedSize:   MaxCursorSize,
		},
		{
			name:           "negative size clamps to default",
			values:         url.Values{"size": {"-5"}},
			expectedCursor: "",
			expectedSize:   DefaultCursorSize,
		},
		{
			name:           "non-numeric size",
			values:         url.Values{"size": {"abc"}},
			expectedCursor: "",
			expectedSize:   DefaultCursorSize,
		},
		{
			name:           "with sort",
			values:         url.Values{"sort": {"created_at,desc", "id,asc"}},
			expectedCursor: "",
			expectedSize:   DefaultCursorSize,
			expectedSort: []Sort{
				{Field: "created_at", Direction: DESC},
				{Field: "id", Direction: ASC},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := CursorRequestFromQuery(tt.values)
			if req.Cursor != tt.expectedCursor {
				t.Errorf("Cursor = %q, want %q", req.Cursor, tt.expectedCursor)
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

func TestNewCursorRequest(t *testing.T) {
	tests := []struct {
		name         string
		cursor       string
		size         int
		expectedSize int
	}{
		{"valid", "abc", 25, 25},
		{"negative size", "abc", -1, DefaultCursorSize},
		{"zero size", "abc", 0, DefaultCursorSize},
		{"size over max", "abc", 5000, MaxCursorSize},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := NewCursorRequest(tt.cursor, tt.size, nil)
			if req.Size != tt.expectedSize {
				t.Errorf("Size = %d, want %d", req.Size, tt.expectedSize)
			}
		})
	}
}

func TestCursorRequestHasCursor(t *testing.T) {
	withCursor := CursorRequest{Cursor: "abc123"}
	if !withCursor.HasCursor() {
		t.Error("HasCursor() should be true")
	}

	withoutCursor := CursorRequest{Cursor: ""}
	if withoutCursor.HasCursor() {
		t.Error("HasCursor() should be false")
	}
}

func TestCursorRequestDecodedCursor(t *testing.T) {
	// Empty cursor returns empty CursorData
	req := CursorRequest{Cursor: ""}
	data, err := req.DecodedCursor()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if data.Value != "" {
		t.Errorf("expected empty value, got %q", data.Value)
	}

	// Valid cursor data
	original := CursorData{Value: "42", Extra: map[string]string{"ts": "12345"}}
	encoded, _ := EncodeCursor(original)
	req = CursorRequest{Cursor: encoded}
	data, err = req.DecodedCursor()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if data.Value != "42" {
		t.Errorf("Value = %q, want %q", data.Value, "42")
	}
	if data.Extra["ts"] != "12345" {
		t.Errorf("Extra[ts] = %q, want %q", data.Extra["ts"], "12345")
	}
}

func TestCursorRequestSortableFields(t *testing.T) {
	// Removes disallowed fields
	req := CursorRequest{
		Size: 20,
		Sort: []Sort{{Field: "name", Direction: ASC}, {Field: "secret", Direction: DESC}},
	}
	req = req.SortableFields("name", "id")
	if len(req.Sort) != 1 || req.Sort[0].Field != "name" {
		t.Errorf("Sort = %v, want [{name asc}]", req.Sort)
	}

	// All disallowed
	req = CursorRequest{
		Size: 20,
		Sort: []Sort{{Field: "secret", Direction: ASC}},
	}
	req = req.SortableFields("name", "id")
	if req.Sort != nil {
		t.Errorf("Sort = %v, want nil", req.Sort)
	}
}

func TestCursorRequestWithDefaultSort(t *testing.T) {
	// Applies default when no sort is set
	req := CursorRequest{Size: 20}
	req = req.WithDefaultSort(Sort{Field: "id", Direction: ASC})
	if len(req.Sort) != 1 || req.Sort[0].Field != "id" {
		t.Errorf("Sort = %v, want [{id asc}]", req.Sort)
	}

	// Does not override existing sort
	req = CursorRequest{Size: 20, Sort: []Sort{{Field: "name", Direction: DESC}}}
	req = req.WithDefaultSort(Sort{Field: "id", Direction: ASC})
	if len(req.Sort) != 1 || req.Sort[0].Field != "name" {
		t.Errorf("Sort = %v, want [{name desc}]", req.Sort)
	}
}

func TestCursorRequestOrderBy(t *testing.T) {
	req := CursorRequest{
		Size: 20,
		Sort: []Sort{
			{Field: "created_at", Direction: DESC},
			{Field: "id", Direction: ASC},
		},
	}
	expected := "created_at desc, id asc"
	if got := req.OrderBy(); got != expected {
		t.Errorf("OrderBy() = %q, want %q", got, expected)
	}

	// No sorts
	req = CursorRequest{Size: 20}
	if got := req.OrderBy(); got != "" {
		t.Errorf("OrderBy() = %q, want empty", got)
	}
}

func TestCursorRequestLimit(t *testing.T) {
	req := CursorRequest{Size: 25}
	if got := req.Limit(); got != 26 {
		t.Errorf("Limit() = %d, want 26", got)
	}
}
