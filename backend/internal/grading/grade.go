package grading

import (
	"errors"
	"fmt"
	"math/big"
)

var (
	// ErrInvalidKey means the answer key itself is broken: a generator bug,
	// never the child's fault.
	ErrInvalidKey = errors.New("invalid answer key")

	// ErrNoSuchOption means the child typed an option number that is not on
	// screen.
	ErrNoSuchOption = errors.New("no such option")
)

// Outcome is the verdict on one answer.
//
// The zero value is OutcomeInvalidInput, so a Result that was never filled
// in can never read as correct.
type Outcome int

const (
	OutcomeInvalidInput  Outcome = iota // could not be read; Result.Err says why
	OutcomeWrong                        // a valid answer with the wrong value
	OutcomeNotSimplified                // the right value, but the exercise wants it simplified
	OutcomeCorrect
)

// String returns the name used in tests and logs. Messages for the child
// are written by the caller, in Portuguese.
func (o Outcome) String() string {
	switch o {
	case OutcomeInvalidInput:
		return "invalid input"
	case OutcomeWrong:
		return "wrong"
	case OutcomeNotSimplified:
		return "not simplified"
	case OutcomeCorrect:
		return "correct"
	default:
		return fmt.Sprintf("Outcome(%d)", int(o))
	}
}

// Key is the answer key of one exercise.
type Key struct {
	Kind Kind

	// Value is the canonical answer, for every kind except KindChoice.
	Value *big.Rat

	// RequireIrreducible makes an equivalent but unsimplified fraction
	// (6/8 for 3/4) OutcomeNotSimplified instead of OutcomeCorrect.
	RequireIrreducible bool

	// Options and Choice are for KindChoice only. Choice is the index of the
	// right option in Options, counted from 0.
	Options []string
	Choice  int
}

// Result is the outcome of grading one answer.
type Result struct {
	Outcome Outcome

	// Expected is the right answer formatted for display, e.g. "3/4" or
	// "0,75", or the text of the right option.
	Expected string

	// Err is why the input could not be read when Outcome is
	// OutcomeInvalidInput (test it with errors.Is), and nil otherwise.
	Err error
}

// Grade grades an answer typed by the child against key.
//
// For KindChoice the child types the option number shown on screen,
// counted from 1.
//
// Anything the child types gives a Result. An error means the key is
// broken (ErrInvalidKey) and the exercise should not have been shown.
func Grade(in string, key Key) (Result, error) {
	expected, err := key.expected()
	if err != nil {
		return Result{}, fmt.Errorf("grade: %w", err)
	}

	if key.Kind == KindChoice {
		return gradeChoice(in, key, expected), nil
	}

	a, err := Parse(in, key.Kind)
	if err != nil {
		return Result{Outcome: OutcomeInvalidInput, Expected: expected, Err: err}, nil
	}

	r := Result{Expected: expected}
	switch {
	case a.Value.Cmp(key.Value) != 0:
		r.Outcome = OutcomeWrong
	case key.RequireIrreducible && !a.InLowestTerms:
		r.Outcome = OutcomeNotSimplified
	default:
		r.Outcome = OutcomeCorrect
	}
	return r, nil
}

func gradeChoice(in string, key Key, expected string) Result {
	a, err := Parse(in, KindInteger)
	if err != nil {
		return Result{Outcome: OutcomeInvalidInput, Expected: expected, Err: err}
	}

	// Parse allows at most 9 digits, so the number fits in an int64.
	n := a.Value.Num().Int64()
	if n < 1 || n > int64(len(key.Options)) {
		err := fmt.Errorf("option %d of %d: %w", n, len(key.Options), ErrNoSuchOption)
		return Result{Outcome: OutcomeInvalidInput, Expected: expected, Err: err}
	}

	if int(n-1) == key.Choice {
		return Result{Outcome: OutcomeCorrect, Expected: expected}
	}
	return Result{Outcome: OutcomeWrong, Expected: expected}
}

// expected checks the key and returns the right answer formatted for
// display.
func (k Key) expected() (string, error) {
	switch k.Kind {
	case KindChoice:
		if k.Choice < 0 || k.Choice >= len(k.Options) {
			return "", fmt.Errorf("%w: choice %d of %d options", ErrInvalidKey, k.Choice, len(k.Options))
		}
		return k.Options[k.Choice], nil

	case KindInteger, KindFraction, KindDecimal:
		s, err := Format(k.Value, k.Kind)
		if err != nil {
			return "", fmt.Errorf("%w: %w", ErrInvalidKey, err)
		}
		// The child must be able to type the right answer: this rejects
		// negative answers and answers past the input limits.
		if _, err := Parse(s, k.Kind); err != nil {
			return "", fmt.Errorf("%w: answer cannot be typed: %w", ErrInvalidKey, err)
		}
		return s, nil

	default:
		return "", fmt.Errorf("%w: %v", ErrInvalidKey, k.Kind)
	}
}
