package exact_test

import (
	"errors"
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

var (
	whichTemplate  = regexp.MustCompile(`^Qual destas frações é equivalente a (\d+)/(\d+)\?$`)
	optionTemplate = regexp.MustCompile(`^(\d+)/(\d+)$`)
)

// Level 3: "Qual destas frações é equivalente a a/b?" with 4 options. The
// oracle rebuilds every option from its misconception rule.
func TestPropertyEquivalentFractionsLevel3(t *testing.T) {
	dens, mults, positions := map[int64]bool{}, map[int64]bool{}, map[int64]bool{}
	oneTermNumerator, oneTermDenominator := 0, 0

	forEachSeed(t, exact.SkillEquivalentFractions, 3, func(e exact.Exercise) {
		n := promptNumbers(t, e, whichTemplate)
		a, b := n[0], n[1]

		if e.Kind != grading.KindChoice || e.Answer != nil || e.RequireIrreducible {
			t.Fatalf("seed %d: kind %v, answer %v, irreducible %v; want choice, nil, false",
				e.Seed, e.Kind, e.Answer, e.RequireIrreducible)
		}
		assertInRange(t, e, "denominator", b, 2, 10)
		if a < 1 || a >= b || oracleGCD(a, b) != 1 {
			t.Fatalf("seed %d: %q: base %d/%d is not proper and irreducible", e.Seed, e.Prompt, a, b)
		}
		if len(e.Options) != 4 {
			t.Fatalf("seed %d: %d options, want 4", e.Seed, len(e.Options))
		}

		// Read every option and find the correct one first: the
		// distractor rules need its multiplier.
		var p, q [4]int64
		var m int64
		correct := 0
		for i, o := range e.Options {
			pq := optionNumbers(t, e, o.Text)
			p[i], q[i] = pq[0], pq[1]
			if !o.Correct {
				continue
			}
			correct++
			positions[int64(i)] = true
			if o.Misconception != "" {
				t.Fatalf("seed %d: correct option %q has misconception %q", e.Seed, o.Text, o.Misconception)
			}
			if q[i]%b != 0 || p[i] != a*(q[i]/b) {
				t.Fatalf("seed %d: correct option %q is not %d/%d scaled", e.Seed, o.Text, a, b)
			}
			m = q[i] / b
		}
		if correct != 1 {
			t.Fatalf("seed %d: %d correct options, want 1", e.Seed, correct)
		}
		assertInRange(t, e, "multiplier", m, 2, 5)

		seen := map[exact.Misconception]int{}
		for i, o := range e.Options {
			if o.Correct {
				continue
			}
			seen[o.Misconception]++

			if p[i]*b == a*q[i] {
				t.Fatalf("seed %d: distractor %q is equivalent to %d/%d", e.Seed, o.Text, a, b)
			}

			switch o.Misconception {
			case exact.MisconceptionAddSame:
				if p[i] != a+m || q[i] != b+m {
					t.Fatalf("seed %d: %q is not %d/%d with %d added to both terms", e.Seed, o.Text, a, b, m)
				}
			case exact.MisconceptionOneTerm:
				switch {
				case p[i] == a*m && q[i] == b:
					oneTermNumerator++
				case p[i] == a && q[i] == b*m:
					oneTermDenominator++
				default:
					t.Fatalf("seed %d: %q is not %d/%d with one term times %d", e.Seed, o.Text, a, b, m)
				}
			case exact.MisconceptionSwap:
				if p[i] != b || q[i] != a {
					t.Fatalf("seed %d: %q is not %d/%d swapped", e.Seed, o.Text, a, b)
				}
			default:
				t.Fatalf("seed %d: distractor %q has misconception %q", e.Seed, o.Text, o.Misconception)
			}
		}
		for _, mc := range []exact.Misconception{exact.MisconceptionAddSame, exact.MisconceptionOneTerm, exact.MisconceptionSwap} {
			if seen[mc] != 1 {
				t.Fatalf("seed %d: misconception %q used %d times, want 1", e.Seed, mc, seen[mc])
			}
		}

		// All four options have different values, so exactly one is right
		// and no two look like the same answer.
		for i := range 4 {
			for j := i + 1; j < 4; j++ {
				if p[i]*q[j] == p[j]*q[i] {
					t.Fatalf("seed %d: options %q and %q have the same value", e.Seed, e.Options[i].Text, e.Options[j].Text)
				}
			}
		}

		// Choosing any distractor is wrong.
		for i, o := range e.Options {
			if o.Correct {
				continue
			}
			res, err := grading.Grade(fmt.Sprint(i+1), e.Key())
			if err != nil || res.Outcome != grading.OutcomeWrong {
				t.Fatalf("seed %d: choosing %q graded %v, %v; want wrong", e.Seed, o.Text, res.Outcome, err)
			}
		}

		dens[b], mults[m] = true, true
	})

	assertCovers(t, "denominator", dens, 2, 10)
	assertCovers(t, "multiplier", mults, 2, 5)
	assertCovers(t, "correct position", positions, 0, 3)
	if oneTermNumerator == 0 || oneTermDenominator == 0 {
		t.Errorf("one-term distractors: %d numerator-only, %d denominator-only; want both", oneTermNumerator, oneTermDenominator)
	}
}

func optionNumbers(t *testing.T, e exact.Exercise, text string) []int64 {
	t.Helper()
	return promptNumbers(t, exact.Exercise{Seed: e.Seed, Prompt: text}, optionTemplate)
}

// A choice exercise must have exactly one correct option, or the grader
// refuses its key instead of guessing.
func TestChoiceKeyNeedsExactlyOneCorrect(t *testing.T) {
	tests := map[string][]exact.Option{
		"none correct": {{Text: "1/2"}, {Text: "2/4"}},
		"two correct":  {{Text: "1/2", Correct: true}, {Text: "2/4", Correct: true}},
		"no options":   nil,
	}
	for name, opts := range tests {
		t.Run(name, func(t *testing.T) {
			e := exact.Exercise{Kind: grading.KindChoice, Options: opts}
			if _, err := grading.Grade("1", e.Key()); !errors.Is(err, grading.ErrInvalidKey) {
				t.Errorf("Grade error = %v, want %v", err, grading.ErrInvalidKey)
			}
		})
	}
}
