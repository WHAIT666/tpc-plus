// Package exact generates maths exercises whose answers are computed, not
// guessed: mode A of TPC+. Every generator is deterministic, so a skill ID,
// generator version, seed and level reproduce an exercise exactly.
package exact

import (
	"errors"
	"math/big"
	"math/rand/v2"

	"github.com/WHAIT666/tpc-plus/backend/internal/grading"
)

// ErrUnknownLevel is returned by Generate for a level the skill does not
// have.
var ErrUnknownLevel = errors.New("unknown level")

// Exercise is one generated exercise with its answer key.
type Exercise struct {
	SkillID string `json:"skill_id"`
	Version int    `json:"version"`
	Seed    uint64 `json:"seed"`
	Level   int    `json:"level"`

	// Prompt is shown to the child, in European Portuguese.
	Prompt string `json:"prompt"`

	Kind grading.Kind `json:"kind"`

	// Answer is the canonical answer for every kind except KindChoice.
	Answer             *big.Rat `json:"answer,omitempty"`
	RequireIrreducible bool     `json:"require_irreducible,omitempty"`

	// Options are the choices of a KindChoice exercise, in display order.
	// Exactly one is correct.
	Options []Option `json:"options,omitempty"`
}

// Option is one choice of a KindChoice exercise.
type Option struct {
	Text    string `json:"text"`
	Correct bool   `json:"correct,omitempty"`

	// Misconception names the mistake that leads a child to this wrong
	// option, so feedback can explain it. Empty for the correct option.
	Misconception Misconception `json:"misconception,omitempty"`
}

// Misconception is a known reasoning mistake that a distractor is built
// from. Distractors come from real mistakes, never random numbers.
type Misconception string

const (
	// MisconceptionAddSame adds the same number to both terms: 3/5 → 5/7.
	MisconceptionAddSame Misconception = "add_same_to_both_terms"

	// MisconceptionOneTerm multiplies only one term: 3/5 → 6/5 or 3/10.
	MisconceptionOneTerm Misconception = "multiply_one_term"

	// MisconceptionSwap swaps numerator and denominator: 3/5 → 5/3.
	MisconceptionSwap Misconception = "swap_terms"
)

// Key returns the answer key the grader checks the child's answer against.
// A choice exercise without exactly one correct option gives a key the
// grader rejects with grading.ErrInvalidKey.
func (e Exercise) Key() grading.Key {
	k := grading.Key{
		Kind:               e.Kind,
		Value:              e.Answer,
		RequireIrreducible: e.RequireIrreducible,
	}
	if e.Kind != grading.KindChoice {
		return k
	}

	k.Choice = -1
	correct := 0
	for i, o := range e.Options {
		k.Options = append(k.Options, o.Text)
		if o.Correct {
			k.Choice = i
			correct++
		}
	}
	if correct != 1 {
		k.Choice = -1
	}
	return k
}

// Generator produces the exercises of one skill.
type Generator interface {
	// SkillID is the catalogue ID, e.g. "mat.fracoes_equivalentes".
	SkillID() string

	// Version is bumped whenever the exercise for some seed and level
	// changes, so stored attempts can still be reproduced.
	Version() int

	// Levels is the number of levels; levels are numbered from 1.
	Levels() int

	// Generate returns the exercise for seed and level. The same version,
	// seed and level give the same exercise on any machine.
	Generate(seed uint64, level int) (Exercise, error)
}

// newRand returns the random source for one exercise. It is always an
// explicitly seeded PCG, never the global functions of math/rand/v2, whose
// seed is random. Mixing in the level makes the same seed give unrelated
// exercises at different levels.
//
// Go does not promise that every helper of rand.Rand returns the same
// values forever; the golden tests catch it if a Go upgrade changes them.
func newRand(seed uint64, level int) *rand.Rand {
	return rand.New(rand.NewPCG(seed, uint64(level)))
}

// between returns a uniform integer in [lo, hi].
func between(r *rand.Rand, lo, hi int) int {
	return lo + r.IntN(hi-lo+1)
}

// pick returns a uniform element of xs, which must not be empty.
func pick[T any](r *rand.Rand, xs []T) T {
	return xs[r.IntN(len(xs))]
}

// coprimes returns the numerators 1..den-1 that make an irreducible
// fraction with den, in increasing order.
func coprimes(den int) []int {
	var out []int
	for n := 1; n < den; n++ {
		if gcd(n, den) == 1 {
			out = append(out, n)
		}
	}
	return out
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}
