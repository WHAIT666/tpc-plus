package grading

import (
	"errors"
	"fmt"
)

// ErrUnknownKind is returned when text names no Kind.
var ErrUnknownKind = errors.New("unknown answer kind")

// Kind is the form an exercise expects its answer in.
//
// The zero value is KindInvalid so that an exercise built without a kind
// fails loudly instead of being graded as if it were an integer.
type Kind int

const (
	KindInvalid  Kind = iota
	KindInteger       // a whole number, e.g. 8
	KindFraction      // a fraction or a whole number, e.g. 3/4 or 2
	KindDecimal       // a decimal or a whole number, e.g. 0,75 or 2
	KindChoice        // one of the options, typed as its number
)

// kinds are the valid kinds, in order.
var kinds = []Kind{KindInteger, KindFraction, KindDecimal, KindChoice}

// String returns the name used in specs, JSON and error messages.
func (k Kind) String() string {
	switch k {
	case KindInteger:
		return "integer"
	case KindFraction:
		return "fraction"
	case KindDecimal:
		return "decimal"
	case KindChoice:
		return "choice"
	default:
		return fmt.Sprintf("Kind(%d)", int(k))
	}
}

// MarshalText writes the kind by name, so JSON (golden files, later the
// API) reads "fraction" instead of 2. An invalid kind is an error rather
// than a silent "Kind(0)".
func (k Kind) MarshalText() ([]byte, error) {
	for _, valid := range kinds {
		if k == valid {
			return []byte(k.String()), nil
		}
	}
	return nil, fmt.Errorf("marshal %v: %w", k, ErrUnknownKind)
}

// UnmarshalText reads a kind written by MarshalText.
func (k *Kind) UnmarshalText(text []byte) error {
	for _, valid := range kinds {
		if string(text) == valid.String() {
			*k = valid
			return nil
		}
	}
	return fmt.Errorf("unmarshal kind %q: %w", text, ErrUnknownKind)
}
