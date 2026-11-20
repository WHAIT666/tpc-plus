package grading

import "fmt"

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
	KindChoice        // the index of one of the options
)

// String returns the name used in specs and error messages.
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
