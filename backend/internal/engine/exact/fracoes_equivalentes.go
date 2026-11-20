package exact

import (
	"fmt"
	"math/big"

	"github.com/WHAIT666/tpc-plus/backend/internal/grading"
)

// SkillEquivalentFractions is the catalogue ID of equivalent fractions.
const SkillEquivalentFractions = "mat.fracoes_equivalentes"

// equivalentFractions generates mat.fracoes_equivalentes.
//
//	Level 1: complete 2/3 = ?/12 (integer answer)
//	Level 2: simplify 18/24 to its irreducible form
//
// Both levels start from a proper, irreducible base fraction and scale it
// by a whole-number multiplier, which is how the textbook builds them.
type equivalentFractions struct{}

func (equivalentFractions) SkillID() string { return SkillEquivalentFractions }
func (equivalentFractions) Version() int    { return 1 }
func (equivalentFractions) Levels() int     { return 2 }

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

	default:
		return Exercise{}, fmt.Errorf("generate %s level %d: %w", g.SkillID(), level, ErrUnknownLevel)
	}

	return e, nil
}
