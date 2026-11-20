package exact

import (
	"fmt"
	"math/big"

	"github.com/WHAIT666/tpc-plus/backend/internal/grading"
)

// SkillAddSubtractFractions is the catalogue ID of adding and subtracting
// fractions.
const SkillAddSubtractFractions = "mat.fracoes_somar_subtrair"

// minus is the typographic minus sign used in textbooks, not a hyphen.
const minus = "−"

// addSubtractFractions generates mat.fracoes_somar_subtrair.
//
//	Level 1: same denominator, 2/7 + 3/7
//	Level 2: one denominator is a multiple of the other, 1/4 + 3/8
//	Level 3: different denominators, 2/3 − 1/4, irreducible answer
//
// Operands are always proper fractions and the result is always positive.
type addSubtractFractions struct{}

func (addSubtractFractions) SkillID() string { return SkillAddSubtractFractions }
func (addSubtractFractions) Version() int    { return 1 }
func (addSubtractFractions) Levels() int     { return 3 }

// multiplePairs are the denominator pairs for level 2: the larger is 2–4
// times the smaller and at most 12.
var multiplePairs = denominatorPairs(func(small, large int) bool {
	return large%small == 0 && large/small <= 4
})

// unrelatedPairs are the denominator pairs for level 3: neither divides the
// other and the least common multiple is at most 24.
var unrelatedPairs = denominatorPairs(func(small, large int) bool {
	return large%small != 0 && lcm(small, large) <= 24
})

func (g addSubtractFractions) Generate(seed uint64, level int) (Exercise, error) {
	r := newRand(seed, level)
	e := Exercise{SkillID: g.SkillID(), Version: g.Version(), Seed: seed, Level: level, Kind: grading.KindFraction}

	var x, y [2]int // the two operands, as numerator and denominator
	switch level {
	case 1:
		// Denominator 3–12. Both operands are proper, so the sum is below 2.
		d := between(r, 3, 12)
		x = [2]int{between(r, 1, d-1), d}
		y = [2]int{between(r, 1, d-1), d}

	case 2:
		p := pick(r, multiplePairs)
		x = [2]int{between(r, 1, p[0]-1), p[0]}
		y = [2]int{between(r, 1, p[1]-1), p[1]}

	case 3:
		p := pick(r, unrelatedPairs)
		x = [2]int{pick(r, coprimes(p[0])), p[0]}
		y = [2]int{pick(r, coprimes(p[1])), p[1]}
		e.RequireIrreducible = true

	default:
		return Exercise{}, fmt.Errorf("generate %s level %d: %w", g.SkillID(), level, ErrUnknownLevel)
	}

	if r.IntN(2) == 1 {
		x, y = y, x
	}
	op := "+"
	if r.IntN(2) == 1 {
		op = minus
	}

	vx, vy := big.NewRat(int64(x[0]), int64(x[1])), big.NewRat(int64(y[0]), int64(y[1]))
	switch c := vx.Cmp(vy); {
	case op == minus && c == 0:
		// 1/2 − 2/4 would be 0: add instead.
		op = "+"
	case op == minus && c < 0:
		// Never negative in 5th grade: subtract the smaller from the larger.
		x, y, vx, vy = y, x, vy, vx
	}

	if op == "+" {
		e.Answer = new(big.Rat).Add(vx, vy)
	} else {
		e.Answer = new(big.Rat).Sub(vx, vy)
	}

	e.Prompt = fmt.Sprintf("Calcula: %d/%d %s %d/%d", x[0], x[1], op, y[0], y[1])
	if e.RequireIrreducible {
		e.Prompt += ". Apresenta o resultado na forma de fração irredutível."
	}
	return e, nil
}

// denominatorPairs returns every (small, large) with 2 ≤ small < large ≤ 12
// that keep accepts, in increasing order.
func denominatorPairs(keep func(small, large int) bool) [][2]int {
	var out [][2]int
	for small := 2; small <= 12; small++ {
		for large := small + 1; large <= 12; large++ {
			if keep(small, large) {
				out = append(out, [2]int{small, large})
			}
		}
	}
	return out
}

func lcm(a, b int) int {
	return a / gcd(a, b) * b
}
