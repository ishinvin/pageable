package pageable

import (
	"testing"
)

func TestParseSort(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *Sort
	}{
		{
			name:     "empty string",
			input:    "",
			expected: nil,
		},
		{
			name:     "whitespace only",
			input:    "   ",
			expected: nil,
		},
		{
			name:     "field only defaults to asc",
			input:    "name",
			expected: &Sort{Field: "name", Direction: ASC},
		},
		{
			name:     "field with asc",
			input:    "name,asc",
			expected: &Sort{Field: "name", Direction: ASC},
		},
		{
			name:     "field with desc",
			input:    "id,desc",
			expected: &Sort{Field: "id", Direction: DESC},
		},
		{
			name:     "case insensitive direction",
			input:    "name,DESC",
			expected: &Sort{Field: "name", Direction: DESC},
		},
		{
			name:     "whitespace handling",
			input:    " name , desc ",
			expected: &Sort{Field: "name", Direction: DESC},
		},
		{
			name:     "invalid direction defaults to asc",
			input:    "name,invalid",
			expected: &Sort{Field: "name", Direction: ASC},
		},
		{
			name:     "sql injection semicolon",
			input:    "id;DROP TABLE users--,asc",
			expected: nil,
		},
		{
			name:     "sql injection quotes",
			input:    "id' OR '1'='1,desc",
			expected: nil,
		},
		{
			name:     "sql injection parentheses",
			input:    "id),desc",
			expected: nil,
		},
		{
			name:     "field with spaces rejected",
			input:    "field name,asc",
			expected: nil,
		},
		{
			name:     "table qualified field allowed",
			input:    "posts.id,desc",
			expected: &Sort{Field: "posts.id", Direction: DESC},
		},
		{
			name:     "underscore field allowed",
			input:    "created_at,asc",
			expected: &Sort{Field: "created_at", Direction: ASC},
		},
		{
			name:     "comma only",
			input:    ",",
			expected: nil,
		},
		{
			name:     "comma with direction only",
			input:    ",desc",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseSort(tt.input)
			if tt.expected == nil {
				if result != nil {
					t.Errorf("expected nil, got %v", result)
				}
				return
			}
			if result == nil {
				t.Fatal("expected non-nil result, got nil")
			}
			if result.Field != tt.expected.Field {
				t.Errorf("Field = %q, want %q", result.Field, tt.expected.Field)
			}
			if result.Direction != tt.expected.Direction {
				t.Errorf("Direction = %q, want %q", result.Direction, tt.expected.Direction)
			}
		})
	}
}

func TestParseSorts(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []Sort
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: nil,
		},
		{
			name:     "empty slice",
			input:    []string{},
			expected: nil,
		},
		{
			name:  "single sort",
			input: []string{"id,desc"},
			expected: []Sort{
				{Field: "id", Direction: DESC},
			},
		},
		{
			name:  "multiple sorts",
			input: []string{"id,desc", "name,asc"},
			expected: []Sort{
				{Field: "id", Direction: DESC},
				{Field: "name", Direction: ASC},
			},
		},
		{
			name:  "field only defaults to asc",
			input: []string{"name"},
			expected: []Sort{
				{Field: "name", Direction: ASC},
			},
		},
		{
			name:     "all empty strings",
			input:    []string{"", "  "},
			expected: nil,
		},
		{
			name:  "mixed valid and empty",
			input: []string{"", "name,desc", "  "},
			expected: []Sort{
				{Field: "name", Direction: DESC},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseSorts(tt.input)
			if tt.expected == nil {
				if result != nil {
					t.Errorf("expected nil, got %v", result)
				}
				return
			}
			if len(result) != len(tt.expected) {
				t.Fatalf("expected %d sorts, got %d: %v", len(tt.expected), len(result), result)
			}
			for i, s := range result {
				if s.Field != tt.expected[i].Field || s.Direction != tt.expected[i].Direction {
					t.Errorf("sort[%d] = %v, want %v", i, s, tt.expected[i])
				}
			}
		})
	}
}

func TestSortString(t *testing.T) {
	tests := []struct {
		sort     Sort
		expected string
	}{
		{Sort{Field: "name", Direction: ASC}, "name,asc"},
		{Sort{Field: "name", Direction: DESC}, "name,desc"},
		{Sort{Field: "id", Direction: ASC}, "id,asc"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.sort.String(); got != tt.expected {
				t.Errorf("Sort.String() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestSortStringRoundTrip(t *testing.T) {
	original := []Sort{
		{Field: "name", Direction: DESC},
		{Field: "id", Direction: ASC},
	}

	// Convert to strings and parse back
	var strs []string
	for _, s := range original {
		strs = append(strs, s.String())
	}

	result := ParseSorts(strs)
	if len(result) != len(original) {
		t.Fatalf("round-trip: expected %d sorts, got %d", len(original), len(result))
	}
	for i, s := range result {
		if s.Field != original[i].Field || s.Direction != original[i].Direction {
			t.Errorf("round-trip: sort[%d] = %v, want %v", i, s, original[i])
		}
	}
}
