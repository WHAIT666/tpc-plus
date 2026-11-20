package grading

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestKindString(t *testing.T) {
	tests := map[Kind]string{
		KindInteger:  "integer",
		KindFraction: "fraction",
		KindDecimal:  "decimal",
		KindChoice:   "choice",
		KindInvalid:  "Kind(0)",
		Kind(42):     "Kind(42)",
	}
	for k, want := range tests {
		if got := k.String(); got != want {
			t.Errorf("Kind(%d).String() = %q, want %q", int(k), got, want)
		}
	}
}

func TestKindJSONRoundTrip(t *testing.T) {
	for _, k := range kinds {
		b, err := json.Marshal(k)
		if err != nil {
			t.Fatalf("Marshal(%v): %v", k, err)
		}
		if want := `"` + k.String() + `"`; string(b) != want {
			t.Errorf("Marshal(%v) = %s, want %s", k, b, want)
		}

		var got Kind
		if err := json.Unmarshal(b, &got); err != nil {
			t.Fatalf("Unmarshal(%s): %v", b, err)
		}
		if got != k {
			t.Errorf("Unmarshal(%s) = %v, want %v", b, got, k)
		}
	}
}

func TestKindJSONRejectsUnknown(t *testing.T) {
	for _, k := range []Kind{KindInvalid, Kind(42)} {
		if _, err := json.Marshal(k); !errors.Is(err, ErrUnknownKind) {
			t.Errorf("Marshal(%v) error = %v, want %v", k, err, ErrUnknownKind)
		}
	}

	for _, in := range []string{`"Kind(0)"`, `"Fraction"`, `""`} {
		var k Kind
		if err := json.Unmarshal([]byte(in), &k); !errors.Is(err, ErrUnknownKind) {
			t.Errorf("Unmarshal(%s) error = %v, want %v", in, err, ErrUnknownKind)
		}
	}
}
