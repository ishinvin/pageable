package pageable

import (
	"encoding/json"
	"testing"
)

type testItem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func TestNewPage(t *testing.T) {
	items := []testItem{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}
	req := PageRequest{Page: 1, Size: 10}
	page := NewPage(items, req, 95)

	if page.Metadata.Page != 1 {
		t.Errorf("Page = %d, want 1", page.Metadata.Page)
	}
	if page.Metadata.Size != 10 {
		t.Errorf("Size = %d, want 10", page.Metadata.Size)
	}
	if page.Metadata.TotalItems != 95 {
		t.Errorf("TotalItems = %d, want 95", page.Metadata.TotalItems)
	}
	if page.Metadata.TotalPages != 10 {
		t.Errorf("TotalPages = %d, want 10", page.Metadata.TotalPages)
	}
	if len(page.Items) != 2 {
		t.Errorf("Items length = %d, want 2", len(page.Items))
	}
}

func TestNewPageTotalPagesCeiling(t *testing.T) {
	tests := []struct {
		name              string
		totalItems        int64
		size              int
		expectedTotalPage int
	}{
		{"exact division", 100, 10, 10},
		{"ceiling needed", 95, 10, 10},
		{"one item", 1, 10, 1},
		{"single page full", 10, 10, 1},
		{"single item over", 11, 10, 2},
		{"zero items", 0, 10, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := PageRequest{Page: 1, Size: tt.size}
			page := NewPage([]testItem{}, req, tt.totalItems)
			if page.Metadata.TotalPages != tt.expectedTotalPage {
				t.Errorf("TotalPages = %d, want %d", page.Metadata.TotalPages, tt.expectedTotalPage)
			}
		})
	}
}

func TestNewPageNilContent(t *testing.T) {
	req := PageRequest{Page: 1, Size: 10}
	page := NewPage[testItem](nil, req, 0)

	if page.Items == nil {
		t.Error("Items should be empty slice, not nil")
	}
	if len(page.Items) != 0 {
		t.Errorf("Items length = %d, want 0", len(page.Items))
	}

	// Verify JSON serialization produces [] not null
	b, err := json.Marshal(page)
	if err != nil {
		t.Fatalf("json.Marshal error: %v", err)
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		t.Fatalf("json.Unmarshal error: %v", err)
	}
	if string(raw["items"]) != "[]" {
		t.Errorf("items JSON = %s, want []", raw["items"])
	}
}

func TestPageJSONSerialization(t *testing.T) {
	req := PageRequest{Page: 1, Size: 2}
	items := []testItem{{ID: 1, Name: "Alice"}}
	page := NewPage(items, req, 5)

	b, err := json.Marshal(page)
	if err != nil {
		t.Fatalf("json.Marshal error: %v", err)
	}

	var result map[string]json.RawMessage
	if err := json.Unmarshal(b, &result); err != nil {
		t.Fatalf("json.Unmarshal error: %v", err)
	}

	// Top-level should have "items" and "metadata"
	if _, ok := result["items"]; !ok {
		t.Error("missing key \"items\" in JSON output")
	}
	if _, ok := result["metadata"]; !ok {
		t.Error("missing key \"metadata\" in JSON output")
	}

	// Check metadata keys use camelCase
	var metadata map[string]any
	if err := json.Unmarshal(result["metadata"], &metadata); err != nil {
		t.Fatalf("json.Unmarshal metadata error: %v", err)
	}
	expectedKeys := []string{"page", "size", "totalItems", "totalPages"}
	for _, key := range expectedKeys {
		if _, ok := metadata[key]; !ok {
			t.Errorf("missing key %q in metadata JSON", key)
		}
	}
}
