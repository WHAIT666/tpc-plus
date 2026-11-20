package exact_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/WHAIT666/tpc-plus/backend/internal/engine/exact"
	"github.com/WHAIT666/tpc-plus/backend/internal/grading"
)

var (
	sumTemplate            = regexp.MustCompile(`^Calcula: (\d+)/(\d+) (?:\+|−) (\d+)/(\d+)$`)
	sumIrreducibleTemplate = regexp.MustCompile(`^Calcula: (\d+)/(\d+) (?:\+|−) (\d+)/(\d+)\. Apresenta o resultado na forma de fração irredutível\.$`)
)

// sumCase is one add/subtract exercise read back from its prompt.
type sumCase struct {
	a, b, c, d int64 // a/b op c/d
	minus      bool
}

// checkSum reads the prompt, recomputes the answer with int64 and checks
// what every level shares: proper operands and a positive result.
func checkSum(t *testing.T, e exact.Exercise, template *regexp.Regexp) sumCase {
	t.Helper()

	n := promptNumbers(t, e, template)
	s := sumCase{a: n[0], b: n[1], c: n[2], d: n[3], minus: strings.Contains(e.Prompt, " − ")}

	if e.Kind != grading.KindFraction {
		t.Fatalf("seed %d: kind %v, want fraction", e.Seed, e.Kind)
	}
	if s.a < 1 || s.a >= s.b || s.c < 1 || s.c >= s.d {
		t.Fatalf("seed %d: %q: operands are not proper fractions", e.Seed, e.Prompt)
	}

	num := s.a*s.d + s.c*s.b
	if s.minus {
		num = s.a*s.d - s.c*s.b
	}
	den := s.b * s.d
	if num <= 0 {
		t.Fatalf("seed %d: %q: result %d/%d is not positive", e.Seed, e.Prompt, num, den)
	}
	g := oracleGCD(num, den)
	assertAnswer(t, e, num/g, den/g)
	return s
}

func TestPropertyAddSubtractFractionsLevel1(t *testing.T) {
	dens, ops := map[int64]bool{}, map[bool]bool{}

	forEachSeed(t, exact.SkillAddSubtractFractions, 1, func(e exact.Exercise) {
		s := checkSum(t, e, sumTemplate)

		if s.b != s.d {
			t.Fatalf("seed %d: %q: denominators differ", e.Seed, e.Prompt)
		}
		assertInRange(t, e, "denominator", s.b, 3, 12)
		if e.RequireIrreducible {
			t.Fatalf("seed %d: level 1 requires irreducible", e.Seed)
		}
		if !s.minus && s.a+s.c >= 2*s.b {
			t.Fatalf("seed %d: %q: sum is not below 2", e.Seed, e.Prompt)
		}

		dens[s.b], ops[s.minus] = true, true
	})

	assertCovers(t, "denominator", dens, 3, 12)
	assertBothOps(t, ops)
}

func TestPropertyAddSubtractFractionsLevel2(t *testing.T) {
	want := oraclePairs(func(small, large int64) bool { return large%small == 0 && large/small <= 4 })
	pairs, ops := map[[2]int64]bool{}, map[bool]bool{}

	forEachSeed(t, exact.SkillAddSubtractFractions, 2, func(e exact.Exercise) {
		s := checkSum(t, e, sumTemplate)

		p := [2]int64{min(s.b, s.d), max(s.b, s.d)}
		if !want[p] {
			t.Fatalf("seed %d: %q: denominators %v are not a 2–4× multiple pair up to 12", e.Seed, e.Prompt, p)
		}
		if e.RequireIrreducible {
			t.Fatalf("seed %d: level 2 requires irreducible", e.Seed)
		}

		pairs[p], ops[s.minus] = true, true
	})

	assertCoversPairs(t, pairs, want)
	assertBothOps(t, ops)
}

func TestPropertyAddSubtractFractionsLevel3(t *testing.T) {
	want := oraclePairs(func(small, large int64) bool {
		return large%small != 0 && small*large/oracleGCD(small, large) <= 24
	})
	pairs, ops := map[[2]int64]bool{}, map[bool]bool{}

	forEachSeed(t, exact.SkillAddSubtractFractions, 3, func(e exact.Exercise) {
		s := checkSum(t, e, sumIrreducibleTemplate)

		p := [2]int64{min(s.b, s.d), max(s.b, s.d)}
		if !want[p] {
			t.Fatalf("seed %d: %q: denominators %v divide each other or have lcm > 24", e.Seed, e.Prompt, p)
		}
		if oracleGCD(s.a, s.b) != 1 || oracleGCD(s.c, s.d) != 1 {
			t.Fatalf("seed %d: %q: operands are not irreducible", e.Seed, e.Prompt)
		}
		if !e.RequireIrreducible {
			t.Fatalf("seed %d: level 3 does not require irreducible", e.Seed)
		}

		// The unsimplified result over b·d must be graded "not simplified"
		// whenever it can be simplified.
		num := s.a*s.d + s.c*s.b
		if s.minus {
			num = s.a*s.d - s.c*s.b
		}
		if oracleGCD(num, s.b*s.d) > 1 {
			res, err := grading.Grade(fmt.Sprintf("%d/%d", num, s.b*s.d), e.Key())
			if err != nil || res.Outcome != grading.OutcomeNotSimplified {
				t.Fatalf("seed %d: unsimplified %d/%d graded %v, %v", e.Seed, num, s.b*s.d, res.Outcome, err)
			}
		}

		pairs[p], ops[s.minus] = true, true
	})

	assertCoversPairs(t, pairs, want)
	assertBothOps(t, ops)
}

// oraclePairs lists (small, large) denominators 2–12 that keep accepts.
func oraclePairs(keep func(small, large int64) bool) map[[2]int64]bool {
	out := map[[2]int64]bool{}
	for small := int64(2); small <= 12; small++ {
		for large := small + 1; large <= 12; large++ {
			if keep(small, large) {
				out[[2]int64{small, large}] = true
			}
		}
	}
	return out
}

func assertCoversPairs(t *testing.T, seen, want map[[2]int64]bool) {
	t.Helper()
	for p := range want {
		if !seen[p] {
			t.Errorf("denominator pair %v never generated in %d seeds", p, propertySeeds)
		}
	}
}

func assertBothOps(t *testing.T, ops map[bool]bool) {
	t.Helper()
	if !ops[false] || !ops[true] {
		t.Errorf("operations seen: + %v, − %v; want both", ops[false], ops[true])
	}
}
