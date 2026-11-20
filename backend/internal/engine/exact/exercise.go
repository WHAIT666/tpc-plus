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

	Kind               grading.Kind `json:"kind"`
	Answer             *big.Rat     `json:"answer"`
	RequireIrreducible bool         `json:"require_irreducible,omitempty"`
}

// Key returns the answer key the grader checks the child's answer against.
func (e Exercise) Key() grading.Key {
	return grading.Key{
		Kind:               e.Kind,
		Value:              e.Answer,
		RequireIrreducible: e.RequireIrreducible,
	}
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
