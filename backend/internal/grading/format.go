package grading

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
)

// ErrNotRepresentable is returned by Format when a value cannot be written
// exactly in the requested kind, such as 3/4 as an integer or 1/3 as a
// decimal.
var ErrNotRepresentable = errors.New("value cannot be written exactly in this kind")

// Format writes v the way the child would type it for an exercise of the
// given kind: "8", "3/4" (or "2" for a whole number) and "0,75" with the
// Portuguese decimal comma. Decimals use as few places as the value needs.
//
// For any answer Parse accepts, parsing Format's output gives the same value.
func Format(v *big.Rat, kind Kind) (string, error) {
	if v == nil {
		return "", fmt.Errorf("format: nil value: %w", ErrNotRepresentable)
	}

	switch kind {
	case KindInteger:
		if !v.IsInt() {
			return "", fmt.Errorf("format %s as %v: %w", v.RatString(), kind, ErrNotRepresentable)
		}
		return v.Num().String(), nil
	case KindFraction:
		return v.RatString(), nil
	case KindDecimal:
		places, ok := decimalPlaces(v.Denom())
		if !ok {
			return "", fmt.Errorf("format %s as %v: %w", v.RatString(), kind, ErrNotRepresentable)
		}
		return strings.Replace(v.FloatString(places), ".", ",", 1), nil
	default:
		return "", fmt.Errorf("format: %v: %w", kind, ErrUnsupportedKind)
	}
}

// decimalPlaces returns the smallest n such that den divides 10^n, and false
// when there is none, i.e. when den has a prime factor other than 2 and 5 and
// the decimal never ends.
func decimalPlaces(den *big.Int) (int, bool) {
	d := new(big.Int).Set(den)

	twos := d.TrailingZeroBits()
	d.Rsh(d, twos)

	fives := 0
	five := big.NewInt(5)
	q, r := new(big.Int), new(big.Int)
	for {
		q.QuoRem(d, five, r)
		if r.Sign() != 0 {
			break
		}
		d.Set(q)
		fives++
	}

	return max(int(twos), fives), d.IsInt64() && d.Int64() == 1
}
