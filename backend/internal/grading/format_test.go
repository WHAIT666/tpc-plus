package grading

import (
	"errors"
	"math/big"
	"testing"
)

func TestFormat(t *testing.T) {
	tests := []struct {
		name    string
		v       *big.Rat
		kind    Kind
		want    string
		wantErr error
	}{
		{name: "integer", v: big.NewRat(8, 1), kind: KindInteger, want: "8"},
		{name: "zero", v: new(big.Rat), kind: KindInteger, want: "0"},
		{name: "fraction", v: big.NewRat(3, 4), kind: KindFraction, want: "3/4"},
		{name: "fraction is reduced", v: big.NewRat(6, 8), kind: KindFraction, want: "3/4"},
		{name: "whole number as fraction", v: big.NewRat(4, 2), kind: KindFraction, want: "2"},
		{name: "decimal uses a comma", v: big.NewRat(3, 4), kind: KindDecimal, want: "0,75"},
		{name: "decimal with three places", v: big.NewRat(1, 8), kind: KindDecimal, want: "0,125"},
		{name: "decimal above one", v: big.NewRat(5, 2), kind: KindDecimal, want: "2,5"},
		{name: "decimal with only fives", v: big.NewRat(1, 25), kind: KindDecimal, want: "0,04"},
		{name: "whole number as decimal", v: big.NewRat(2, 1), kind: KindDecimal, want: "2"},
		{name: "nine decimal places", v: big.NewRat(123456789, 1_000_000_000), kind: KindDecimal, want: "0,123456789"},

		{name: "fraction as integer", v: big.NewRat(3, 4), kind: KindInteger, wantErr: ErrNotRepresentable},
		{name: "repeating decimal", v: big.NewRat(1, 3), kind: KindDecimal, wantErr: ErrNotRepresentable},
		{name: "repeating decimal with twos and fives", v: big.NewRat(1, 30), kind: KindDecimal, wantErr: ErrNotRepresentable},
		{name: "nil value", v: nil, kind: KindFraction, wantErr: ErrNotRepresentable},
		{name: "choice", v: big.NewRat(1, 1), kind: KindChoice, wantErr: ErrUnsupportedKind},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Format(tt.v, tt.kind)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("Format(%v, %v) error = %v, want %v", tt.v, tt.kind, err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Format(%v, %v) unexpected error: %v", tt.v, tt.kind, err)
			}
			if got != tt.want {
				t.Errorf("Format(%v, %v) = %q, want %q", tt.v, tt.kind, got, tt.want)
			}
		})
	}
}

// The fuzz test in step 4 checks this for any input; these are the cases
// worth reading.
func TestFormatRoundTrip(t *testing.T) {
	tests := []struct {
		in   string
		kind Kind
	}{
		{"8", KindInteger},
		{"6/8", KindFraction},
		{"0/5", KindFraction},
		{"999999999/1", KindFraction},
		{"0,50", KindDecimal},
		{"0.05", KindDecimal},
		{"999999999,999999999", KindDecimal},
	}

	for _, tt := range tests {
		a, err := Parse(tt.in, tt.kind)
		if err != nil {
			t.Fatalf("Parse(%q, %v): %v", tt.in, tt.kind, err)
		}
		s, err := Format(a.Value, tt.kind)
		if err != nil {
			t.Fatalf("Format(%v, %v): %v", a.Value, tt.kind, err)
		}
		b, err := Parse(s, tt.kind)
		if err != nil {
			t.Fatalf("Parse(Format(%q)) = Parse(%q): %v", tt.in, s, err)
		}
		if a.Value.Cmp(b.Value) != 0 {
			t.Errorf("round trip of %q via %q: got %v, want %v", tt.in, s, b.Value, a.Value)
		}
	}
}
