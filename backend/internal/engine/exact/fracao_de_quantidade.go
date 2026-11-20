package exact

import (
	"fmt"
	"math/big"

	"github.com/WHAIT666/tpc-plus/backend/internal/grading"
)

// SkillFractionOfQuantity is the catalogue ID of a fraction of a quantity.
const SkillFractionOfQuantity = "mat.fracao_de_quantidade"

// fractionOfQuantity generates mat.fracao_de_quantidade.
//
//	Level 1: unit fraction of a quantity, 1/4 de 20
//	Level 2: any fraction of a quantity, 3/5 de 35
//	Level 3: a short word problem
//
// The quantity is always a multiple of the denominator, so every answer is
// a whole number.
type fractionOfQuantity struct{}

func (fractionOfQuantity) SkillID() string { return SkillFractionOfQuantity }
func (fractionOfQuantity) Version() int    { return 1 }
func (fractionOfQuantity) Levels() int     { return 3 }

// wordProblem is a level 3 template. Quantities stay realistic for the
// story: a class has 20–30 pupils, not 90.
type wordProblem struct {
	// text takes the name (if any), the quantity and the fraction, in the
	// order given by format.
	format   func(name string, q, num, den int) string
	min, max int

	// remaining asks for what is left (q − part) instead of the part.
	remaining bool
}

// wordProblems are the level 3 templates. Order matters: the generator
// picks by index, so reordering changes the output for a seed.
var wordProblems = []wordProblem{
	{
		format: func(name string, q, num, den int) string {
			return fmt.Sprintf("%s tinha %d cromos e deu %d/%d ao irmão. Com quantos cromos ficou?", name, q, num, den)
		},
		min: 12, max: 60, remaining: true,
	},
	{
		format: func(_ string, q, num, den int) string {
			return fmt.Sprintf("Numa turma de %d alunos, %d/%d vão de autocarro para a escola. Quantos alunos vão de autocarro?", q, num, den)
		},
		min: 20, max: 30,
	},
	{
		format: func(name string, q, num, den int) string {
			return fmt.Sprintf("%s já leu %d/%d de um livro com %d páginas. Quantas páginas faltam ler?", name, num, den, q)
		},
		min: 40, max: 120, remaining: true,
	},
	{
		format: func(name string, q, num, den int) string {
			return fmt.Sprintf("%s comprou %d rebuçados e comeu %d/%d. Quantos rebuçados comeu?", name, q, num, den)
		},
		min: 10, max: 40,
	},
}

// names come with their article, as Portuguese uses them: "A Rita tinha".
// They are invented for the exercises and are not anyone real.
var names = []string{
	"A Rita", "O Tomás", "A Beatriz", "O Duarte",
	"A Matilde", "O Rodrigo", "A Leonor", "O Afonso",
}

func (g fractionOfQuantity) Generate(seed uint64, level int) (Exercise, error) {
	r := newRand(seed, level)
	e := Exercise{SkillID: g.SkillID(), Version: g.Version(), Seed: seed, Level: level, Kind: grading.KindInteger}

	switch level {
	case 1:
		// Unit fractions 1/2–1/10, quantity 2–10 times the denominator.
		den := between(r, 2, 10)
		q := den * between(r, 2, 10)

		e.Prompt = fmt.Sprintf("Calcula: 1/%d de %d", den, q)
		e.Answer = big.NewRat(int64(q/den), 1)

	case 2:
		// Non-unit irreducible fractions, denominator 3–10 (2 has none).
		den := between(r, 3, 10)
		num := pick(r, coprimes(den)[1:])
		q := den * between(r, 2, 10)

		e.Prompt = fmt.Sprintf("Calcula: %d/%d de %d", num, den, q)
		e.Answer = big.NewRat(int64(q/den*num), 1)

	case 3:
		wp := pick(r, wordProblems)
		name := pick(r, names)
		den := between(r, 2, 10)
		num := pick(r, coprimes(den))
		// The smallest and largest multiples of den inside the story's range.
		q := den * between(r, (wp.min+den-1)/den, wp.max/den)

		part := q / den * num
		e.Prompt = wp.format(name, q, num, den)
		if wp.remaining {
			e.Answer = big.NewRat(int64(q-part), 1)
		} else {
			e.Answer = big.NewRat(int64(part), 1)
		}

	default:
		return Exercise{}, fmt.Errorf("generate %s level %d: %w", g.SkillID(), level, ErrUnknownLevel)
	}

	return e, nil
}
