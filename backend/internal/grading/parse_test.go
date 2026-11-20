package grading

import (
	"errors"
	"strings"
	"testing"
)

// parseTests is shared with FuzzParse, which uses every input as a seed.
var parseTests = []struct {
	name       string
	in         string
	kind       Kind
	want       string // exact value as big.Rat.RatString, when no error
	wantLowest bool
	wantErr    error
}{
	// Rows of the parsing table in docs/m1-spec.md.
	{name: "fraction", in: "3/4", kind: KindFraction, want: "3/4", wantLowest: true},
	{name: "fraction with spaces", in: " 3 / 4 ", kind: KindFraction, want: "3/4", wantLowest: true},
	{name: "unsimplified fraction", in: "6/8", kind: KindFraction, want: "3/4", wantLowest: false},
	{name: "integer as fraction", in: "2", kind: KindFraction, want: "2", wantLowest: true},
	{name: "decimal comma", in: "0,75", kind: KindDecimal, want: "3/4", wantLowest: true},
	{name: "decimal point", in: "0.75", kind: KindDecimal, want: "3/4", wantLowest: true},
	{name: "zero denominator", in: "1/0", kind: KindFraction, wantErr: ErrZeroDenominator},
	{name: "zero denominator on integer exercise", in: "1/0", kind: KindInteger, wantErr: ErrZeroDenominator},
	{name: "negative fraction", in: "-3/4", kind: KindFraction, wantErr: ErrNegative},
	{name: "negative integer", in: "-3", kind: KindInteger, wantErr: ErrNegative},
	{name: "negative decimal", in: "-0,5", kind: KindDecimal, wantErr: ErrNegative},
	{name: "negative denominator", in: "3/-4", kind: KindFraction, wantErr: ErrNegative},
	{name: "mixed number", in: "1 1/2", kind: KindFraction, wantErr: ErrMixedNumber},
	{name: "mixed number on decimal exercise", in: "1 1/2", kind: KindDecimal, wantErr: ErrMixedNumber},
	{name: "thousands separator and comma", in: "1.000,5", kind: KindDecimal, wantErr: ErrMalformed},
	{name: "double slash", in: "3//4", kind: KindFraction, wantErr: ErrMalformed},
	{name: "letters", in: "abc", kind: KindInteger, wantErr: ErrMalformed},
	{name: "empty", in: "", kind: KindInteger, wantErr: ErrEmpty},
	{name: "only spaces", in: "   ", kind: KindInteger, wantErr: ErrEmpty},

	// Input limits: 32 bytes and 9 digits per number.
	{name: "32 bytes", in: "1/2" + strings.Repeat(" ", 29), kind: KindFraction, want: "1/2", wantLowest: true},
	{name: "33 bytes", in: "1/2" + strings.Repeat(" ", 30), kind: KindFraction, wantErr: ErrTooLong},
	{name: "9 digits", in: "999999999", kind: KindInteger, want: "999999999", wantLowest: true},
	{name: "10 digits", in: "1000000000", kind: KindInteger, wantErr: ErrTooManyDigits},
	{name: "9 digit numerator and denominator", in: "123456789/987654321", kind: KindFraction, want: "13717421/109739369", wantLowest: false},
	{name: "10 digit denominator", in: "1/1000000000", kind: KindFraction, wantErr: ErrTooManyDigits},
	{name: "9 decimal places", in: "0,123456789", kind: KindDecimal, want: "123456789/1000000000", wantLowest: true},
	{name: "10 decimal places", in: "0,1234567891", kind: KindDecimal, wantErr: ErrTooManyDigits},

	// Leading zeros: only a lone 0 may start a number; decimal places are free.
	{name: "zero", in: "0", kind: KindInteger, want: "0", wantLowest: true},
	{name: "leading zero integer", in: "07", kind: KindInteger, wantErr: ErrLeadingZero},
	{name: "leading zero numerator", in: "010/3", kind: KindFraction, wantErr: ErrLeadingZero},
	{name: "leading zero denominator", in: "3/04", kind: KindFraction, wantErr: ErrLeadingZero},
	{name: "leading zero decimal", in: "00,5", kind: KindDecimal, wantErr: ErrLeadingZero},
	{name: "zeros after the decimal mark", in: "0,05", kind: KindDecimal, want: "1/20", wantLowest: true},
	{name: "zero numerator", in: "0/5", kind: KindFraction, want: "0", wantLowest: false},

	// Forms big.Rat.SetString accepts but a 5th grader never types.
	{name: "plus sign", in: "+3/4", kind: KindFraction, wantErr: ErrMalformed},
	{name: "exponent", in: "1e5", kind: KindInteger, wantErr: ErrMalformed},
	{name: "hex prefix", in: "0x10", kind: KindInteger, wantErr: ErrMalformed},
	{name: "underscore", in: "1_000", kind: KindInteger, wantErr: ErrMalformed},
	{name: "no digits before the mark", in: ".5", kind: KindDecimal, wantErr: ErrMalformed},
	{name: "no digits after the mark", in: "5.", kind: KindDecimal, wantErr: ErrMalformed},
	{name: "full-width digits", in: "３/４", kind: KindFraction, wantErr: ErrMalformed},
	{name: "unicode minus", in: "−3", kind: KindInteger, wantErr: ErrMalformed},

	// Other malformed shapes.
	{name: "missing denominator", in: "3/", kind: KindFraction, wantErr: ErrMalformed},
	{name: "missing numerator", in: "/4", kind: KindFraction, wantErr: ErrMalformed},
	{name: "two slashes", in: "1/2/3", kind: KindFraction, wantErr: ErrMalformed},
	{name: "two decimal marks", in: "1,2,3", kind: KindDecimal, wantErr: ErrMalformed},
	{name: "space inside a decimal", in: "0 ,75", kind: KindDecimal, wantErr: ErrMalformed},
	{name: "space inside a number", in: "1 2", kind: KindInteger, wantErr: ErrMalformed},
	{name: "decimal fraction", in: "0,5/2", kind: KindFraction, wantErr: ErrMalformed},

	// Whitespace the keyboard or terminal adds.
	{name: "windows line ending", in: "3/4\r\n", kind: KindFraction, want: "3/4", wantLowest: true},
	{name: "tab around the answer", in: "\t8\t", kind: KindInteger, want: "8", wantLowest: true},
	{name: "space before slash only", in: "3 /4", kind: KindFraction, want: "3/4", wantLowest: true},

	// The typed form must match the exercise kind.
	{name: "integer", in: "8", kind: KindInteger, want: "8", wantLowest: true},
	{name: "integer as decimal", in: "2", kind: KindDecimal, want: "2", wantLowest: true},
	{name: "fraction on integer exercise", in: "8/1", kind: KindInteger, wantErr: ErrWrongForm},
	{name: "decimal on integer exercise", in: "8,0", kind: KindInteger, wantErr: ErrWrongForm},
	{name: "decimal on fraction exercise", in: "0,75", kind: KindFraction, wantErr: ErrWrongForm},
	{name: "fraction on decimal exercise", in: "3/4", kind: KindDecimal, wantErr: ErrWrongForm},

	// Kinds the parser does not handle.
	{name: "choice", in: "1", kind: KindChoice, wantErr: ErrUnsupportedKind},
	{name: "zero kind", in: "1", kind: KindInvalid, wantErr: ErrUnsupportedKind},
}

func TestParse(t *testing.T) {
	for _, tt := range parseTests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.in, tt.kind)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("Parse(%q, %v) error = %v, want %v", tt.in, tt.kind, err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Parse(%q, %v) unexpected error: %v", tt.in, tt.kind, err)
			}
			if s := got.Value.RatString(); s != tt.want {
				t.Errorf("Parse(%q, %v) value = %s, want %s", tt.in, tt.kind, s, tt.want)
			}
			if got.InLowestTerms != tt.wantLowest {
				t.Errorf("Parse(%q, %v) InLowestTerms = %v, want %v", tt.in, tt.kind, got.InLowestTerms, tt.wantLowest)
			}
		})
	}
}

func TestParseErrorDoesNotEchoLongInput(t *testing.T) {
	in := strings.Repeat("9", 10_000)

	_, err := Parse(in, KindInteger)

	if !errors.Is(err, ErrTooLong) {
		t.Fatalf("error = %v, want %v", err, ErrTooLong)
	}
	if strings.Contains(err.Error(), in[:100]) {
		t.Errorf("error message echoes the input: %.80s...", err.Error())
	}
}

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
