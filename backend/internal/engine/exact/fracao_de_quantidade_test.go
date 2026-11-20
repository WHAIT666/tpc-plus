package exact_test

import (
	"regexp"
	"testing"

	"github.com/WHAIT666/tpc-plus/backend/internal/engine/exact"
	"github.com/WHAIT666/tpc-plus/backend/internal/grading"
)

var (
	unitOfTemplate     = regexp.MustCompile(`^Calcula: 1/(\d+) de (\d+)$`)
	fractionOfTemplate = regexp.MustCompile(`^Calcula: (\d+)/(\d+) de (\d+)$`)
)

// story is the oracle's own copy of one word problem: how to read the
// numbers back and what the question asks.
type story struct {
	re                *regexp.Regexp
	qAt, numAt, denAt int // positions of quantity and fraction in the match
	min, max          int64
	remaining         bool
}

// The name is "A" or "O" plus one word; the oracle does not care which.
var stories = []story{
	{
		re:  regexp.MustCompile(`^(?:A|O) \pL+ tinha (\d+) cromos e deu (\d+)/(\d+) ao irmão\. Com quantos cromos ficou\?$`),
		qAt: 0, numAt: 1, denAt: 2, min: 12, max: 60, remaining: true,
	},
	{
		re:  regexp.MustCompile(`^Numa turma de (\d+) alunos, (\d+)/(\d+) vão de autocarro para a escola\. Quantos alunos vão de autocarro\?$`),
		qAt: 0, numAt: 1, denAt: 2, min: 20, max: 30,
	},
	{
		re:  regexp.MustCompile(`^(?:A|O) \pL+ já leu (\d+)/(\d+) de um livro com (\d+) páginas\. Quantas páginas faltam ler\?$`),
		qAt: 2, numAt: 0, denAt: 1, min: 40, max: 120, remaining: true,
	},
	{
		re:  regexp.MustCompile(`^(?:A|O) \pL+ comprou (\d+) rebuçados e comeu (\d+)/(\d+)\. Quantos rebuçados comeu\?$`),
		qAt: 0, numAt: 1, denAt: 2, min: 10, max: 40,
	},
}

// partOf checks q is a whole multiple of den and returns num/den of q.
func partOf(t *testing.T, e exact.Exercise, num, den, q int64) int64 {
	t.Helper()
	if e.Kind != grading.KindInteger || e.RequireIrreducible {
		t.Fatalf("seed %d: kind %v, irreducible %v; want integer, false", e.Seed, e.Kind, e.RequireIrreducible)
	}
	if q%den != 0 {
		t.Fatalf("seed %d: %q: %d is not a multiple of %d, answer would not be whole", e.Seed, e.Prompt, q, den)
	}
	return q / den * num
}

func TestPropertyFractionOfQuantityLevel1(t *testing.T) {
	dens := map[int64]bool{}

	forEachSeed(t, exact.SkillFractionOfQuantity, 1, func(e exact.Exercise) {
		n := promptNumbers(t, e, unitOfTemplate)
		den, q := n[0], n[1]

		assertInRange(t, e, "denominator", den, 2, 10)
		part := partOf(t, e, 1, den, q)
		assertInRange(t, e, "quantity ÷ denominator", q/den, 2, 10)
		assertAnswer(t, e, part, 1)

		dens[den] = true
	})

	assertCovers(t, "denominator", dens, 2, 10)
}

func TestPropertyFractionOfQuantityLevel2(t *testing.T) {
	dens := map[int64]bool{}

	forEachSeed(t, exact.SkillFractionOfQuantity, 2, func(e exact.Exercise) {
		n := promptNumbers(t, e, fractionOfTemplate)
		num, den, q := n[0], n[1], n[2]

		assertInRange(t, e, "denominator", den, 3, 10)
		if num < 2 || num >= den || oracleGCD(num, den) != 1 {
			t.Fatalf("seed %d: %q: %d/%d is not a proper, irreducible, non-unit fraction", e.Seed, e.Prompt, num, den)
		}
		part := partOf(t, e, num, den, q)
		assertInRange(t, e, "quantity ÷ denominator", q/den, 2, 10)
		assertAnswer(t, e, part, 1)

		dens[den] = true
	})

	assertCovers(t, "denominator", dens, 3, 10)
}

func TestPropertyFractionOfQuantityLevel3(t *testing.T) {
	dens, used := map[int64]bool{}, map[int64]bool{}

	forEachSeed(t, exact.SkillFractionOfQuantity, 3, func(e exact.Exercise) {
		matched := -1
		for i, s := range stories {
			if s.re.MatchString(e.Prompt) {
				matched = i
			}
		}
		if matched < 0 {
			t.Fatalf("seed %d: prompt %q matches no known story", e.Seed, e.Prompt)
		}
		s := stories[matched]

		n := promptNumbers(t, e, s.re)
		q, num, den := n[s.qAt], n[s.numAt], n[s.denAt]

		assertInRange(t, e, "denominator", den, 2, 10)
		if num < 1 || num >= den || oracleGCD(num, den) != 1 {
			t.Fatalf("seed %d: %q: %d/%d is not proper and irreducible", e.Seed, e.Prompt, num, den)
		}
		assertInRange(t, e, "quantity", q, s.min, s.max)

		want := partOf(t, e, num, den, q)
		if s.remaining {
			want = q - want
		}
		if want <= 0 {
			t.Fatalf("seed %d: %q: answer %d is not positive", e.Seed, e.Prompt, want)
		}
		assertAnswer(t, e, want, 1)

		dens[den], used[int64(matched)] = true, true
	})

	assertCovers(t, "denominator", dens, 2, 10)
	assertCovers(t, "story", used, 0, int64(len(stories)-1))
}
