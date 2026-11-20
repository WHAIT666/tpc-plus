package exact

import (
	"fmt"
	"math/big"
	"math/rand/v2"

	"github.com/WHAIT666/tpc-plus/backend/internal/grading"
)

// SkillEquivalentFractions is the catalogue ID of equivalent fractions.
const SkillEquivalentFractions = "mat.fracoes_equivalentes"

// equivalentFractions generates mat.fracoes_equivalentes.
//
//	Level 1: complete 2/3 = ?/12 (integer answer)
//	Level 2: simplify 18/24 to its irreducible form
//	Level 3: which fraction is equivalent to 3/5? (4 options)
//
// Every level starts from a proper, irreducible base fraction and scales it
// by a whole-number multiplier, which is how the textbook builds them.
type equivalentFractions struct{}

func (equivalentFractions) SkillID() string { return SkillEquivalentFractions }

// Version 1 covers levels 1–3. Adding level 3 did not change the output of
// any existing seed and level, so it did not need a bump.
func (equivalentFractions) Version() int { return 1 }
func (equivalentFractions) Levels() int  { return 3 }

func (g equivalentFractions) Generate(seed uint64, level int) (Exercise, error) {
	r := newRand(seed, level)
	e := Exercise{SkillID: g.SkillID(), Version: g.Version(), Seed: seed, Level: level}

	switch level {
	case 1:
		// Denominator 2–10, multiplier 2–5.
		den := between(r, 2, 10)
		num := pick(r, coprimes(den))
		m := between(r, 2, 5)

		e.Prompt = fmt.Sprintf("Completa: %d/%d = ?/%d", num, den, den*m)
		e.Kind = grading.KindInteger
		e.Answer = big.NewRat(int64(num*m), 1)

	case 2:
		// Irreducible base with denominator 2–12, multiplier 2–6.
		den := between(r, 2, 12)
		num := pick(r, coprimes(den))
		m := between(r, 2, 6)

		e.Prompt = fmt.Sprintf("Simplifica %d/%d. Escreve a fração irredutível.", num*m, den*m)
		e.Kind = grading.KindFraction
		e.Answer = big.NewRat(int64(num), int64(den))
		e.RequireIrreducible = true

	case 3:
		// Same ranges as level 1: denominator 2–10, multiplier 2–5.
		den := between(r, 2, 10)
		num := pick(r, coprimes(den))
		m := between(r, 2, 5)

		e.Prompt = fmt.Sprintf("Qual destas frações é equivalente a %d/%d?", num, den)
		e.Kind = grading.KindChoice
		e.Options = equivalentOptions(r, num, den, m)

	default:
		return Exercise{}, fmt.Errorf("generate %s level %d: %w", g.SkillID(), level, ErrUnknownLevel)
	}

	return e, nil
}

// equivalentOptions returns num·m/den·m and one distractor per
// misconception, shuffled.
//
// None of the distractors is equivalent to num/den, because the base is
// proper and m ≥ 2. They also have different values from each other, with
// one exception handled here: multiplying only the numerator can equal
// adding m to both terms (1/4, m=2: 2/4 and 3/6) or swapping the terms
// (1/2, m=4: 4/2 and 2/1). Multiplying only the denominator never collides,
// so it is used instead.
func equivalentOptions(r *rand.Rand, num, den, m int) []Option {
	add := [2]int{num + m, den + m}
	swap := [2]int{den, num}

	one := [2]int{num * m, den}
	if r.IntN(2) == 1 || sameValue(one, add) || sameValue(one, swap) {
		one = [2]int{num, den * m}
	}

	opts := []Option{
		{Text: fraction(num*m, den*m), Correct: true},
		{Text: fraction(add[0], add[1]), Misconception: MisconceptionAddSame},
		{Text: fraction(one[0], one[1]), Misconception: MisconceptionOneTerm},
		{Text: fraction(swap[0], swap[1]), Misconception: MisconceptionSwap},
	}
	r.Shuffle(len(opts), func(i, j int) { opts[i], opts[j] = opts[j], opts[i] })
	return opts
}

// fraction writes n/d as the child would see it on paper, unsimplified.
func fraction(n, d int) string {
	return fmt.Sprintf("%d/%d", n, d)
}

func sameValue(a, b [2]int) bool {
	return a[0]*b[1] == b[0]*a[1]
}
