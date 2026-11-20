package grading

import (
	"errors"
	"testing"
)

// parseErrors are every error Parse may wrap. Each one maps to a message for
// the child, so an error outside this list would reach her as a blank.
var parseErrors = []error{
	ErrEmpty, ErrTooLong, ErrTooManyDigits, ErrMalformed, ErrLeadingZero,
	ErrNegative, ErrMixedNumber, ErrZeroDenominator, ErrWrongForm, ErrUnsupportedKind,
}

// FuzzParse checks that Parse never panics on any input or kind, and that
// whatever it accepts is a sane value that survives Format and Parse again.
//
//	go test -run=^$ -fuzz=FuzzParse -fuzztime=60s ./internal/grading
func FuzzParse(f *testing.F) {
	for _, tt := range parseTests {
		f.Add(tt.in, uint8(tt.kind))
	}

	f.Fuzz(func(t *testing.T, in string, k uint8) {
		// Cover every kind, including KindInvalid and one unknown value.
		kind := Kind(k % 6)

		a, err := Parse(in, kind)
		if err != nil {
			if !isParseError(err) {
				t.Fatalf("Parse(%q, %v) error %v wraps none of the known errors", in, kind, err)
			}
			return
		}

		if len(in) > maxInputLen {
			t.Fatalf("Parse accepted %d bytes, limit %d", len(in), maxInputLen)
		}
		if a.Value == nil {
			t.Fatalf("Parse(%q, %v) accepted with a nil value", in, kind)
		}
		if a.Value.Sign() < 0 {
			t.Fatalf("Parse(%q, %v) accepted a negative value %v", in, kind, a.Value)
		}

		s, err := Format(a.Value, kind)
		if err != nil {
			t.Fatalf("Format(%v, %v) of accepted %q: %v", a.Value, kind, in, err)
		}
		b, err := Parse(s, kind)
		if err != nil {
			t.Fatalf("Parse(Format(%q)) = Parse(%q, %v): %v", in, s, kind, err)
		}
		if a.Value.Cmp(b.Value) != 0 {
			t.Fatalf("round trip of %q via %q: got %v, want %v", in, s, b.Value, a.Value)
		}
		if !b.InLowestTerms {
			t.Fatalf("Format(%v, %v) = %q is not in lowest terms", a.Value, kind, s)
		}
	})
}

func isParseError(err error) bool {
	for _, target := range parseErrors {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}
