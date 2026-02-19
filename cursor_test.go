package pageable

import (
	"strings"
	"testing"
)

func TestCursorRoundTrip(t *testing.T) {
	tests := []string{
		"123",
		"abc-def",
		"2024-01-15T10:30:00Z",
		"",
		"some/path/value",
		"unicode-日本語",
	}

	for _, value := range tests {
		t.Run(value, func(t *testing.T) {
			encoded, err := EncodeCursor(CursorData{Value: value})
			if err != nil {
				t.Fatalf("EncodeCursor error: %v", err)
			}
			decoded, err := DecodeCursor(encoded)
			if err != nil {
				t.Fatalf("DecodeCursor error: %v", err)
			}
			if decoded.Value != value {
				t.Errorf("decoded = %q, want %q", decoded.Value, value)
			}
		})
	}
}

func TestEncodeCursorURLSafe(t *testing.T) {
	encoded, _ := EncodeCursor(CursorData{Value: "test-value-123"})

	// base64 URL encoding should not contain + or /
	if strings.ContainsAny(encoded, "+/") {
		t.Errorf("cursor %q contains non-URL-safe characters", encoded)
	}
}

func TestEncodeCursorDeterministic(t *testing.T) {
	a, _ := EncodeCursor(CursorData{Value: "test-123"})
	b, _ := EncodeCursor(CursorData{Value: "test-123"})
	if a != b {
		t.Errorf("non-deterministic encoding: %q != %q", a, b)
	}
}

func TestEncodeCursorWithExtra(t *testing.T) {
	original := CursorData{
		Value: "user-42",
		Extra: map[string]string{
			"created_at": "2024-01-15T10:30:00Z",
			"name":       "alice",
		},
	}

	encoded, err := EncodeCursor(original)
	if err != nil {
		t.Fatalf("EncodeCursor error: %v", err)
	}

	decoded, err := DecodeCursor(encoded)
	if err != nil {
		t.Fatalf("DecodeCursor error: %v", err)
	}

	if decoded.Value != original.Value {
		t.Errorf("Value = %q, want %q", decoded.Value, original.Value)
	}
	if len(decoded.Extra) != len(original.Extra) {
		t.Fatalf("Extra length = %d, want %d", len(decoded.Extra), len(original.Extra))
	}
	for k, v := range original.Extra {
		if decoded.Extra[k] != v {
			t.Errorf("Extra[%q] = %q, want %q", k, decoded.Extra[k], v)
		}
	}
}

func TestEncodeCursorWithoutExtra(t *testing.T) {
	original := CursorData{Value: "simple-value"}

	encoded, err := EncodeCursor(original)
	if err != nil {
		t.Fatalf("EncodeCursor error: %v", err)
	}

	decoded, err := DecodeCursor(encoded)
	if err != nil {
		t.Fatalf("DecodeCursor error: %v", err)
	}

	if decoded.Value != original.Value {
		t.Errorf("Value = %q, want %q", decoded.Value, original.Value)
	}
	if decoded.Extra != nil {
		t.Errorf("Extra = %v, want nil", decoded.Extra)
	}
}

func TestDecodeCursorInvalidBase64(t *testing.T) {
	_, err := DecodeCursor("!!!invalid!!!")
	if err == nil {
		t.Error("expected error for invalid base64")
	}
}

func TestDecodeCursorInvalidJSON(t *testing.T) {
	// Encode raw non-JSON bytes as base64
	_, err := DecodeCursor("bm90LWpzb24=") // "not-json" in base64
	if err == nil {
		t.Error("expected error for invalid JSON inside valid base64")
	}
	if !strings.Contains(err.Error(), "invalid cursor data") {
		t.Errorf("unexpected error message: %v", err)
	}
}
