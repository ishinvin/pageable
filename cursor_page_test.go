package pageable

import (
	"encoding/json"
	"testing"
)

func TestNewCursorPage(t *testing.T) {
	items := []testItem{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}
	nextCursor, _ := EncodeCursor(CursorData{Value: "2"})

	page := NewCursorPage(items, nextCursor, "", true, false, 10)

	if len(page.Items) != 2 {
		t.Errorf("Items length = %d, want 2", len(page.Items))
	}
	if page.Metadata.NextCursor != nextCursor {
		t.Errorf("NextCursor = %q, want %q", page.Metadata.NextCursor, nextCursor)
	}
	if page.Metadata.PrevCursor != "" {
		t.Errorf("PrevCursor = %q, want empty", page.Metadata.PrevCursor)
	}
	if !page.Metadata.HasNext {
		t.Error("HasNext should be true")
	}
	if page.Metadata.HasPrev {
		t.Error("HasPrev should be false")
	}
	if page.Metadata.Size != 10 {
		t.Errorf("Size = %d, want 10", page.Metadata.Size)
	}
}

func TestNewCursorPageNoMore(t *testing.T) {
	items := []testItem{{ID: 1, Name: "Alice"}}
	page := NewCursorPage(items, "", "", false, false, 10)

	if page.Metadata.HasNext {
		t.Error("HasNext should be false")
	}
	if page.Metadata.HasPrev {
		t.Error("HasPrev should be false")
	}
	if page.Metadata.NextCursor != "" {
		t.Errorf("NextCursor = %q, want empty", page.Metadata.NextCursor)
	}
}

func TestNewCursorPageWithPrev(t *testing.T) {
	items := []testItem{{ID: 2, Name: "Bob"}}
	prevCursor, _ := EncodeCursor(CursorData{Value: "1", Direction: Prev})

	page := NewCursorPage(items, "", prevCursor, false, true, 10)

	if page.Metadata.HasNext {
		t.Error("HasNext should be false")
	}
	if !page.Metadata.HasPrev {
		t.Error("HasPrev should be true")
	}
	if page.Metadata.PrevCursor != prevCursor {
		t.Errorf("PrevCursor = %q, want %q", page.Metadata.PrevCursor, prevCursor)
	}
}

func TestNewCursorPageNilContent(t *testing.T) {
	page := NewCursorPage[testItem](nil, "", "", false, false, 10)

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

func TestCursorPageJSONEmptyCursors(t *testing.T) {
	page := NewCursorPage([]testItem{{ID: 1}}, "", "", false, false, 10)

	b, err := json.Marshal(page)
	if err != nil {
		t.Fatalf("json.Marshal error: %v", err)
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		t.Fatalf("json.Unmarshal error: %v", err)
	}

	var metadata map[string]json.RawMessage
	if err := json.Unmarshal(raw["metadata"], &metadata); err != nil {
		t.Fatalf("json.Unmarshal metadata error: %v", err)
	}

	// Empty cursors should be present as empty strings
	if string(metadata["nextCursor"]) != `""` {
		t.Errorf("nextCursor = %s, want empty string", metadata["nextCursor"])
	}
	if string(metadata["prevCursor"]) != `""` {
		t.Errorf("prevCursor = %s, want empty string", metadata["prevCursor"])
	}
}
