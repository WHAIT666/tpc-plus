// Package grading parses answers typed by a 5th grader and grades them
// exactly against the canonical answer of an exercise.
package grading

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

const (
	maxInputLen = 32 // bytes, before trimming
	maxDigits   = 9  // per run of digits; keeps every value well inside int64
)

// Errors returned by Parse. Callers test for them with errors.Is; each one
// maps to a different message for the child.
var (
	ErrEmpty           = errors.New("empty answer")
	ErrTooLong         = errors.New("answer too long")
	ErrTooManyDigits   = errors.New("too many digits in a number")
	ErrMalformed       = errors.New("not a number")
	ErrLeadingZero     = errors.New("number starts with a zero")
	ErrNegative        = errors.New("negative numbers are not allowed")
	ErrMixedNumber     = errors.New("mixed numbers are not allowed")
	ErrZeroDenominator = errors.New("zero denominator")
	ErrWrongForm       = errors.New("answer form does not match the exercise")
	ErrUnsupportedKind = errors.New("kind is not parsed from text")
)

// Answer is a parsed answer.
type Answer struct {
	// Value is the exact value, always in lowest terms.
	Value *big.Rat

	// InLowestTerms is false only when the child typed a fraction that can
	// still be simplified, such as 6/8. Value has already been reduced, so
	// this is the only place the grader can learn it.
	InLowestTerms bool
}

// Parse reads an answer typed on a maths keypad for an exercise of the given
// kind. It accepts whole numbers, fractions such as "3/4" (spaces around the
// slash are fine) and decimals with a comma or a point ("0,75", "0.75").
//
// It never uses big.Rat.SetString, which would also accept signs, exponents,
// hex and octal prefixes and underscores. Every input is first checked
// against the rules in docs/m1-spec.md and only then converted.
func Parse(in string, kind Kind) (Answer, error) {
	if kind != KindInteger && kind != KindFraction && kind != KindDecimal {
		return Answer{}, fmt.Errorf("parse: %v: %w", kind, ErrUnsupportedKind)
	}
	// Checked before anything else touches the input, and the input is not
	// echoed back: it may be huge.
	if len(in) > maxInputLen {
		return Answer{}, fmt.Errorf("parse: %d bytes, limit %d: %w", len(in), maxInputLen, ErrTooLong)
	}

	s := strings.TrimSpace(in)
	if s == "" {
		return Answer{}, fmt.Errorf("parse: %w", ErrEmpty)
	}

	a, form, err := parseForm(s)
	if err != nil {
		return Answer{}, fmt.Errorf("parse %q: %w", s, err)
	}
	if !accepts(kind, form) {
		return Answer{}, fmt.Errorf("parse %q: typed a %v, want a %v: %w", s, form, kind, ErrWrongForm)
	}
	return a, nil
}

// accepts reports whether an answer typed in the given form is valid for an
// exercise of the given kind. A whole number is a valid fraction and a valid
// decimal; nothing else crosses kinds.
func accepts(kind, form Kind) bool {
	return form == kind || form == KindInteger
}

// parseForm works out which form s was typed in and converts it.
func parseForm(s string) (Answer, Kind, error) {
	if strings.HasPrefix(s, "-") {
		return Answer{}, 0, ErrNegative
	}
	if num, den, ok := strings.Cut(s, "/"); ok {
		a, err := parseFraction(strings.TrimRight(num, " "), strings.TrimLeft(den, " "))
		return a, KindFraction, err
	}
	if i := strings.IndexAny(s, ",."); i >= 0 {
		a, err := parseDecimal(s[:i], s[i+1:])
		return a, KindDecimal, err
	}
	n, err := parseWhole(s)
	if err != nil {
		return Answer{}, 0, err
	}
	return Answer{Value: new(big.Rat).SetInt64(n), InLowestTerms: true}, KindInteger, nil
}

func parseFraction(num, den string) (Answer, error) {
	if strings.HasPrefix(den, "-") {
		return Answer{}, ErrNegative
	}
	if whole, frac, ok := strings.Cut(num, " "); ok && isDigits(whole) && isDigits(strings.TrimLeft(frac, " ")) {
		return Answer{}, ErrMixedNumber
	}

	n, err := parseWhole(num)
	if err != nil {
		return Answer{}, fmt.Errorf("numerator: %w", err)
	}
	d, err := parseWhole(den)
	if err != nil {
		return Answer{}, fmt.Errorf("denominator: %w", err)
	}
	if d == 0 {
		return Answer{}, ErrZeroDenominator
	}

	v := new(big.Rat).SetFrac64(n, d)
	// SetFrac64 reduces, so the typed fraction was irreducible exactly when
	// the denominator survived unchanged.
	return Answer{Value: v, InLowestTerms: v.Denom().IsInt64() && v.Denom().Int64() == d}, nil
}

func parseDecimal(whole, frac string) (Answer, error) {
	if _, err := parseWhole(whole); err != nil {
		return Answer{}, fmt.Errorf("whole part: %w", err)
	}
	// Zeros right after the decimal mark are normal ("0,05"), so the
	// leading-zero rule does not apply here.
	if err := checkDigits(frac); err != nil {
		return Answer{}, fmt.Errorf("decimal places: %w", err)
	}

	// At most 9 + 9 digits, which fits in an int64.
	n, err := strconv.ParseInt(whole+frac, 10, 64)
	if err != nil {
		return Answer{}, fmt.Errorf("decimal: %w", err)
	}
	d := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(len(frac))), nil)
	v := new(big.Rat).SetFrac(big.NewInt(n), d)
	return Answer{Value: v, InLowestTerms: true}, nil
}

// parseWhole converts a whole number: ASCII digits only, at most maxDigits,
// and no leading zero unless the number is 0 itself.
func parseWhole(s string) (int64, error) {
	if err := checkDigits(s); err != nil {
		return 0, err
	}
	if len(s) > 1 && s[0] == '0' {
		return 0, ErrLeadingZero
	}
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("whole number: %w", err)
	}
	return n, nil
}

// checkDigits checks that s is a non-empty run of at most maxDigits ASCII
// digits.
func checkDigits(s string) error {
	if !isDigits(s) {
		return ErrMalformed
	}
	if len(s) > maxDigits {
		return ErrTooManyDigits
	}
	return nil
}

// isDigits reports whether s is non-empty and made only of ASCII digits.
// unicode.IsDigit is not used on purpose: it accepts "３" and "٣".
func isDigits(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}
