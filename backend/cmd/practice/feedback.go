package main

import (
	"errors"
	"fmt"

	"github.com/WHAIT666/tpc-plus/backend/internal/engine/exact"
	"github.com/WHAIT666/tpc-plus/backend/internal/grading"
)

// Everything in this file is shown to a 10-year-old: European Portuguese,
// short sentences, and always what to do next rather than what went wrong.

// misconceptionHints explain the mistake behind each distractor.
var misconceptionHints = map[exact.Misconception]string{
	exact.MisconceptionAddSame: "Somar o mesmo número aos dois termos não dá uma fração equivalente. " +
		"Multiplica ou divide os dois termos pelo mesmo número.",
	exact.MisconceptionOneTerm: "Tens de multiplicar os dois termos pelo mesmo número, não só um.",
	exact.MisconceptionSwap:    "Trocaste o numerador com o denominador.",
}

// feedback is the message after a graded answer.
func feedback(e exact.Exercise, typed string, res grading.Result) string {
	switch res.Outcome {
	case grading.OutcomeCorrect:
		return "Certo!"
	case grading.OutcomeNotSimplified:
		return fmt.Sprintf("Quase! Está certo, mas falta simplificar. A fração irredutível é %s.", res.Expected)
	}

	if e.Kind != grading.KindChoice {
		return fmt.Sprintf("Ainda não. A resposta certa é %s.", res.Expected)
	}

	msg := fmt.Sprintf("Ainda não. A resposta certa é a opção %d: %s.", correctOption(e), res.Expected)
	if hint := chosenMisconceptionHint(e, typed); hint != "" {
		msg += "\n" + hint
	}
	return msg
}

// correctOption returns the number shown next to the correct option.
func correctOption(e exact.Exercise) int {
	for i, o := range e.Options {
		if o.Correct {
			return i + 1
		}
	}
	return 0
}

// chosenMisconceptionHint explains the mistake behind the option the child
// chose, or returns "" when there is nothing to explain.
func chosenMisconceptionHint(e exact.Exercise, typed string) string {
	a, err := grading.Parse(typed, grading.KindInteger)
	if err != nil {
		return ""
	}
	n := a.Value.Num().Int64()
	if n < 1 || n > int64(len(e.Options)) {
		return ""
	}
	return misconceptionHints[e.Options[n-1].Misconception]
}

// invalidHint tells the child how to type the answer after the parser
// rejected it.
func invalidHint(e exact.Exercise, err error) string {
	switch {
	case errors.Is(err, grading.ErrEmpty):
		return "Escreve uma resposta."
	case errors.Is(err, grading.ErrNoSuchOption):
		return fmt.Sprintf("Escolhe o número de uma das opções, de 1 a %d.", len(e.Options))
	case errors.Is(err, grading.ErrNegative):
		return "Aqui não há números negativos."
	case errors.Is(err, grading.ErrZeroDenominator):
		return "O denominador não pode ser 0."
	case errors.Is(err, grading.ErrMixedNumber):
		return "Escreve a resposta só como fração, por exemplo 3/2."
	case errors.Is(err, grading.ErrLeadingZero):
		return "Escreve o número sem zeros à esquerda."
	case errors.Is(err, grading.ErrTooLong), errors.Is(err, grading.ErrTooManyDigits):
		return "Esse número é grande demais."
	default:
		return "Não percebi. " + formHint(e.Kind)
	}
}

// formHint shows the shape of answer the exercise expects.
func formHint(k grading.Kind) string {
	switch k {
	case grading.KindInteger:
		return "Escreve só um número inteiro, por exemplo 8."
	case grading.KindFraction:
		return "Escreve uma fração, por exemplo 3/4."
	case grading.KindDecimal:
		return "Escreve um número decimal, por exemplo 0,75."
	case grading.KindChoice:
		return "Escreve o número da opção."
	default:
		return ""
	}
}
