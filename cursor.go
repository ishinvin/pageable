package pageable

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// CursorDirection indicates forward or backward traversal.
type CursorDirection string

const (
	// Next indicates forward pagination (items after the cursor).
	Next CursorDirection = "next"
	// Prev indicates backward pagination (items before the cursor).
	Prev CursorDirection = "prev"
)

// CursorData holds the raw values encoded inside a cursor token.
type CursorData struct {
	// Value is the primary cursor value (e.g., an ID or timestamp).
	Value string `json:"v"`
	// Direction indicates whether this cursor paginates forward or backward.
	Direction CursorDirection `json:"d"`
	// Extra holds additional cursor fields for compound cursors
	// (e.g., created_at + id for stable ordering).
	Extra map[string]string `json:"e"`
}

// EncodeCursor encodes a CursorData struct into a base64 URL-safe cursor token.
func EncodeCursor(data CursorData) (string, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("pageable: failed to encode cursor data: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// DecodeCursor decodes a base64 cursor token back to CursorData.
func DecodeCursor(cursor string) (CursorData, error) {
	b, err := base64.URLEncoding.DecodeString(cursor)
	if err != nil {
		return CursorData{}, fmt.Errorf("pageable: invalid cursor encoding: %w", err)
	}
	var data CursorData
	if err := json.Unmarshal(b, &data); err != nil {
		return CursorData{}, fmt.Errorf("pageable: invalid cursor data: %w", err)
	}
	return data, nil
}
