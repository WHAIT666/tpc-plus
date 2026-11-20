package exact_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/WHAIT666/tpc-plus/backend/internal/engine/exact"
	"github.com/WHAIT666/tpc-plus/backend/internal/grading"
)

var (
	completeTemplate = regexp.MustCompile(`^Completa: (\d+)/(\d+) = \?/(\d+)$`)
	simplifyTemplate = regexp.MustCompile(`^Simplifica (\d+)/(\d+)\. Escreve a fração irredutível\.$`)
)

// Level 1: "Completa: a/b = ?/c". The answer is a·c/b, which must be whole.
func TestPropertyEquivalentFractionsLevel1(t *testing.T) {
	dens, mults := map[int64]bool{}, map[int64]bool{}

	forEachSeed(t, exact.SkillEquivalentFractions, 1, func(e exact.Exercise) {
		n := promptNumbers(t, e, completeTemplate)
		a, b, c := n[0], n[1], n[2]

		if e.Kind != grading.KindInteger || e.RequireIrreducible {
			t.Fatalf("seed %d: kind %v, irreducible %v; want integer, false", e.Seed, e.Kind, e.RequireIrreducible)
		}
		assertInRange(t, e, "denominator", b, 2, 10)
		if a < 1 || a >= b || oracleGCD(a, b) != 1 {
			t.Fatalf("seed %d: %q: base %d/%d is not proper and irreducible", e.Seed, e.Prompt, a, b)
		}
		if c%b != 0 {
			t.Fatalf("seed %d: %q: %d is not a multiple of %d", e.Seed, e.Prompt, c, b)
		}
		m := c / b
		assertInRange(t, e, "multiplier", m, 2, 5)
		if !ratIsInt(e.Answer) {
			t.Fatalf("seed %d: answer %v is not whole", e.Seed, e.Answer)
		}
		assertAnswer(t, e, a*c/b, 1)

		dens[b], mults[m] = true, true
	})

	assertCovers(t, "denominator", dens, 2, 10)
	assertCovers(t, "multiplier", mults, 2, 5)
}

// Level 2: "Simplifica p/q". The answer is p/q in lowest terms, and the
// shown fraction must really need simplifying.
func TestPropertyEquivalentFractionsLevel2(t *testing.T) {
	dens, mults := map[int64]bool{}, map[int64]bool{}

	forEachSeed(t, exact.SkillEquivalentFractions, 2, func(e exact.Exercise) {
		n := promptNumbers(t, e, simplifyTemplate)
		p, q := n[0], n[1]

		if e.Kind != grading.KindFraction || !e.RequireIrreducible {
			t.Fatalf("seed %d: kind %v, irreducible %v; want fraction, true", e.Seed, e.Kind, e.RequireIrreducible)
		}

		// The base is irreducible, so the common factor is the multiplier.
		g := oracleGCD(p, q)
		num, den := p/g, q/g
		assertInRange(t, e, "multiplier", g, 2, 6)
		assertInRange(t, e, "denominator", den, 2, 12)
		if num < 1 || num >= den {
			t.Fatalf("seed %d: %q: base %d/%d is not proper", e.Seed, e.Prompt, num, den)
		}
		assertAnswer(t, e, num, den)

		// Typing the shown fraction back must not count as simplified.
		res, err := grading.Grade(fmt.Sprintf("%d/%d", p, q), e.Key())
		if err != nil {
			t.Fatalf("seed %d: grade unsimplified: %v", e.Seed, err)
		}
		if res.Outcome != grading.OutcomeNotSimplified {
			t.Fatalf("seed %d: typing %d/%d back graded %v, want %v", e.Seed, p, q, res.Outcome, grading.OutcomeNotSimplified)
		}

		dens[den], mults[g] = true, true
	})

	assertCovers(t, "denominator", dens, 2, 12)
	assertCovers(t, "multiplier", mults, 2, 6)
}
