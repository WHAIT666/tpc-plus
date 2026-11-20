package exact_test

// The oracle lives in package exact_test, not exact, so the compiler stops it
// from reusing the generators' helpers. It recomputes every answer from the
// prompt the child actually sees, with int64 and its own gcd, never big.Rat.

import (
	"math/big"
	"regexp"
	"strconv"
	"testing"

	"github.com/WHAIT666/tpc-plus/backend/internal/engine/exact"
	"github.com/WHAIT666/tpc-plus/backend/internal/grading"
)

// propertySeeds is how many seeds every skill and level is checked with.
const propertySeeds = 10_000

func oracleGCD(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

// promptNumbers matches prompt against a template and returns the numbers
// captured by its groups.
func promptNumbers(t *testing.T, e exact.Exercise, template *regexp.Regexp) []int64 {
	t.Helper()

	m := template.FindStringSubmatch(e.Prompt)
	if m == nil {
		t.Fatalf("seed %d: prompt %q does not match %s", e.Seed, e.Prompt, template)
	}
	nums := make([]int64, 0, len(m)-1)
	for _, s := range m[1:] {
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			t.Fatalf("seed %d: number %q in prompt: %v", e.Seed, s, err)
		}
		nums = append(nums, n)
	}
	return nums
}

// assertAnswer checks the exercise's answer against the oracle's num/den,
// which must already be in lowest terms.
func assertAnswer(t *testing.T, e exact.Exercise, num, den int64) {
	t.Helper()

	if e.Answer == nil {
		t.Fatalf("seed %d: nil answer", e.Seed)
	}
	gotNum, gotDen := e.Answer.Num(), e.Answer.Denom()
	if !gotNum.IsInt64() || !gotDen.IsInt64() || gotNum.Int64() != num || gotDen.Int64() != den {
		t.Fatalf("seed %d: %q answer = %v, oracle says %d/%d", e.Seed, e.Prompt, e.Answer, num, den)
	}
}

func assertInRange(t *testing.T, e exact.Exercise, what string, v, lo, hi int64) {
	t.Helper()
	if v < lo || v > hi {
		t.Fatalf("seed %d: %q: %s = %d, want %d–%d", e.Seed, e.Prompt, what, v, lo, hi)
	}
}

// assertCovers checks that every value in lo..hi was generated at least
// once, so a range that is allowed is also reached.
func assertCovers(t *testing.T, what string, seen map[int64]bool, lo, hi int64) {
	t.Helper()
	for v := lo; v <= hi; v++ {
		if !seen[v] {
			t.Errorf("%s = %d never generated in %d seeds", what, v, propertySeeds)
		}
	}
}

// forEachSeed generates seeds 1..propertySeeds at one level, checks what
// every exercise must satisfy, and hands each one to check.
func forEachSeed(t *testing.T, id string, level int, check func(exact.Exercise)) {
	t.Helper()

	reg, err := exact.Default()
	if err != nil {
		t.Fatal(err)
	}
	g, err := reg.Generator(id)
	if err != nil {
		t.Fatal(err)
	}

	for seed := uint64(1); seed <= propertySeeds; seed++ {
		e, err := g.Generate(seed, level)
		if err != nil {
			t.Fatalf("Generate(%d, %d): %v", seed, level, err)
		}
		if e.SkillID != id || e.Version != g.Version() || e.Seed != seed || e.Level != level {
			t.Fatalf("seed %d: identity = %s v%d seed %d level %d", seed, e.SkillID, e.Version, e.Seed, e.Level)
		}
		assertGradesOwnAnswerCorrect(t, e)
		check(e)
	}
}

// assertGradesOwnAnswerCorrect types the canonical answer the way the child
// would and expects the grader to accept it.
func assertGradesOwnAnswerCorrect(t *testing.T, e exact.Exercise) {
	t.Helper()

	typed, err := grading.Format(e.Answer, e.Kind)
	if err != nil {
		t.Fatalf("seed %d: format answer %v as %v: %v", e.Seed, e.Answer, e.Kind, err)
	}
	res, err := grading.Grade(typed, e.Key())
	if err != nil {
		t.Fatalf("seed %d: invalid key: %v", e.Seed, err)
	}
	if res.Outcome != grading.OutcomeCorrect {
		t.Fatalf("seed %d: own answer %q graded %v", e.Seed, typed, res.Outcome)
	}
}

// ratIsInt reports whether r is a whole number, for invariants that only
// need to read the exercise.
func ratIsInt(r *big.Rat) bool { return r != nil && r.IsInt() }
