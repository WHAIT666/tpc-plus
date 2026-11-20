package grading

import (
	"errors"
	"math/big"
	"testing"
)

func TestGrade(t *testing.T) {
	threeQuarters := Key{Kind: KindFraction, Value: big.NewRat(3, 4)}
	threeQuartersIrreducible := Key{Kind: KindFraction, Value: big.NewRat(3, 4), RequireIrreducible: true}
	two := Key{Kind: KindFraction, Value: big.NewRat(2, 1), RequireIrreducible: true}
	eight := Key{Kind: KindInteger, Value: big.NewRat(8, 1)}
	decimal := Key{Kind: KindDecimal, Value: big.NewRat(3, 4)}
	choice := Key{Kind: KindChoice, Options: []string{"5/7", "6/10", "6/5", "5/3"}, Choice: 1}

	tests := []struct {
		name         string
		in           string
		key          Key
		want         Outcome
		wantExpected string
		wantErr      error // Result.Err, only for OutcomeInvalidInput
	}{
		{name: "fraction correct", in: "3/4", key: threeQuarters, want: OutcomeCorrect, wantExpected: "3/4"},
		{name: "equivalent fraction accepted", in: "6/8", key: threeQuarters, want: OutcomeCorrect, wantExpected: "3/4"},
		{name: "fraction wrong", in: "1/2", key: threeQuarters, want: OutcomeWrong, wantExpected: "3/4"},
		{name: "fraction malformed", in: "abc", key: threeQuarters, want: OutcomeInvalidInput, wantExpected: "3/4", wantErr: ErrMalformed},
		{name: "decimal on fraction exercise", in: "0,75", key: threeQuarters, want: OutcomeInvalidInput, wantExpected: "3/4", wantErr: ErrWrongForm},
		{name: "negative on fraction exercise", in: "-3/4", key: threeQuarters, want: OutcomeInvalidInput, wantExpected: "3/4", wantErr: ErrNegative},

		{name: "irreducible correct", in: "3/4", key: threeQuartersIrreducible, want: OutcomeCorrect, wantExpected: "3/4"},
		{name: "not simplified", in: "6/8", key: threeQuartersIrreducible, want: OutcomeNotSimplified, wantExpected: "3/4"},
		{name: "wrong beats not simplified", in: "4/8", key: threeQuartersIrreducible, want: OutcomeWrong, wantExpected: "3/4"},
		{name: "whole number answer", in: "2", key: two, want: OutcomeCorrect, wantExpected: "2"},
		{name: "whole number typed as fraction", in: "2/1", key: two, want: OutcomeCorrect, wantExpected: "2"},
		{name: "whole number not simplified", in: "4/2", key: two, want: OutcomeNotSimplified, wantExpected: "2"},

		{name: "integer correct", in: "8", key: eight, want: OutcomeCorrect, wantExpected: "8"},
		{name: "integer wrong", in: "9", key: eight, want: OutcomeWrong, wantExpected: "8"},
		{name: "integer with leading zero", in: "08", key: eight, want: OutcomeInvalidInput, wantExpected: "8", wantErr: ErrLeadingZero},
		{name: "fraction on integer exercise", in: "16/2", key: eight, want: OutcomeInvalidInput, wantExpected: "8", wantErr: ErrWrongForm},

		{name: "decimal comma", in: "0,75", key: decimal, want: OutcomeCorrect, wantExpected: "0,75"},
		{name: "decimal point", in: "0.75", key: decimal, want: OutcomeCorrect, wantExpected: "0,75"},
		{name: "decimal trailing zero", in: "0,750", key: decimal, want: OutcomeCorrect, wantExpected: "0,75"},
		{name: "decimal wrong", in: "0,7", key: decimal, want: OutcomeWrong, wantExpected: "0,75"},
		{name: "fraction on decimal exercise", in: "3/4", key: decimal, want: OutcomeInvalidInput, wantExpected: "0,75", wantErr: ErrWrongForm},

		// Choices are typed as the option number shown to the child, from 1.
		{name: "choice correct", in: "2", key: choice, want: OutcomeCorrect, wantExpected: "6/10"},
		{name: "choice wrong", in: "1", key: choice, want: OutcomeWrong, wantExpected: "6/10"},
		{name: "choice last option", in: "4", key: choice, want: OutcomeWrong, wantExpected: "6/10"},
		{name: "choice zero", in: "0", key: choice, want: OutcomeInvalidInput, wantExpected: "6/10", wantErr: ErrNoSuchOption},
		{name: "choice past the end", in: "5", key: choice, want: OutcomeInvalidInput, wantExpected: "6/10", wantErr: ErrNoSuchOption},
		{name: "choice malformed", in: "b", key: choice, want: OutcomeInvalidInput, wantExpected: "6/10", wantErr: ErrMalformed},
		{name: "choice empty", in: " ", key: choice, want: OutcomeInvalidInput, wantExpected: "6/10", wantErr: ErrEmpty},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Grade(tt.in, tt.key)
			if err != nil {
				t.Fatalf("Grade(%q) unexpected error: %v", tt.in, err)
			}

			if got.Outcome != tt.want {
				t.Errorf("Grade(%q) outcome = %v, want %v", tt.in, got.Outcome, tt.want)
			}
			if got.Expected != tt.wantExpected {
				t.Errorf("Grade(%q) expected = %q, want %q", tt.in, got.Expected, tt.wantExpected)
			}
			switch {
			case tt.wantErr == nil && got.Err != nil:
				t.Errorf("Grade(%q) Result.Err = %v, want nil", tt.in, got.Err)
			case tt.wantErr != nil && !errors.Is(got.Err, tt.wantErr):
				t.Errorf("Grade(%q) Result.Err = %v, want %v", tt.in, got.Err, tt.wantErr)
			}
		})
	}
}

// A broken key is a bug in a generator, not something the child did, so it
// is an error from Grade rather than an outcome.
func TestGradeInvalidKey(t *testing.T) {
	tests := []struct {
		name string
		key  Key
	}{
		{name: "zero key", key: Key{}},
		{name: "nil value", key: Key{Kind: KindFraction}},
		{name: "fraction value on integer kind", key: Key{Kind: KindInteger, Value: big.NewRat(3, 4)}},
		{name: "repeating decimal", key: Key{Kind: KindDecimal, Value: big.NewRat(1, 3)}},
		{name: "negative answer", key: Key{Kind: KindFraction, Value: big.NewRat(-3, 4)}},
		{name: "answer too long to type", key: Key{Kind: KindInteger, Value: big.NewRat(1_000_000_000, 1)}},
		{name: "no options", key: Key{Kind: KindChoice}},
		{name: "choice past the end", key: Key{Kind: KindChoice, Options: []string{"a", "b"}, Choice: 2}},
		{name: "negative choice", key: Key{Kind: KindChoice, Options: []string{"a", "b"}, Choice: -1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Grade("1", tt.key)
			if !errors.Is(err, ErrInvalidKey) {
				t.Errorf("Grade error = %v, want %v", err, ErrInvalidKey)
			}
		})
	}
}

func TestOutcomeString(t *testing.T) {
	tests := map[Outcome]string{
		OutcomeInvalidInput:  "invalid input",
		OutcomeWrong:         "wrong",
		OutcomeNotSimplified: "not simplified",
		OutcomeCorrect:       "correct",
		Outcome(42):          "Outcome(42)",
	}
	for o, want := range tests {
		if got := o.String(); got != want {
			t.Errorf("Outcome(%d).String() = %q, want %q", int(o), got, want)
		}
	}
}
