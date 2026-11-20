package main

import (
	"bytes"
	"strconv"
	"strings"
	"testing"

	"github.com/WHAIT666/tpc-plus/backend/internal/engine/exact"
	"github.com/WHAIT666/tpc-plus/backend/internal/grading"
)

// exercises returns the exercises a session with these settings shows.
func exercises(t *testing.T, skill string, level, n int, seed uint64) []exact.Exercise {
	t.Helper()
	reg, err := exact.Default()
	if err != nil {
		t.Fatal(err)
	}
	g, err := reg.Generator(skill)
	if err != nil {
		t.Fatal(err)
	}
	var out []exact.Exercise
	for i := range n {
		e, err := g.Generate(seed+uint64(i), level)
		if err != nil {
			t.Fatal(err)
		}
		out = append(out, e)
	}
	return out
}

// rightAnswer types the correct answer the way the child would.
func rightAnswer(t *testing.T, e exact.Exercise) string {
	t.Helper()
	if e.Kind == grading.KindChoice {
		return strconv.Itoa(correctOption(e))
	}
	s, err := grading.Format(e.Answer, e.Kind)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

// play runs the CLI with args and the given lines as keyboard input.
func play(t *testing.T, args []string, lines ...string) string {
	t.Helper()
	var out, errOut bytes.Buffer
	in := strings.NewReader(strings.Join(lines, "\r\n") + "\r\n") // typed on Windows
	if err := run(args, in, &out, &errOut); err != nil {
		t.Fatalf("run(%v): %v\nstderr: %s", args, err, errOut.String())
	}
	return out.String()
}

func TestSessionAllCorrect(t *testing.T) {
	for _, skill := range []string{exact.SkillEquivalentFractions, exact.SkillAddSubtractFractions, exact.SkillFractionOfQuantity} {
		for level := 1; level <= 3; level++ {
			var answers []string
			for _, e := range exercises(t, skill, level, 5, 42) {
				answers = append(answers, rightAnswer(t, e))
			}

			args := []string{"-skill", skill, "-level", strconv.Itoa(level), "-n", "5", "-seed", "42"}
			out := play(t, args, answers...)

			if !strings.Contains(out, "Resultado: 5 de 5 certas.") {
				t.Errorf("%s level %d: want 5 of 5, got:\n%s", skill, level, out)
			}
			if !strings.Contains(out, "-seed 42") {
				t.Errorf("%s level %d: replay line missing:\n%s", skill, level, out)
			}
		}
	}
}

func TestSessionShowsOptions(t *testing.T) {
	e := exercises(t, exact.SkillEquivalentFractions, 3, 1, 7)[0]
	out := play(t, []string{"-level", "3", "-n", "1", "-seed", "7"}, rightAnswer(t, e))

	for i, o := range e.Options {
		if want := "  " + strconv.Itoa(i+1) + ") " + o.Text + "\n"; !strings.Contains(out, want) {
			t.Errorf("option %q missing from:\n%s", want, out)
		}
	}
}

func TestInvalidInputIsNotAnAttempt(t *testing.T) {
	e := exercises(t, exact.SkillEquivalentFractions, 2, 1, 1)[0]
	out := play(t, []string{"-level", "2", "-n", "1", "-seed", "1"}, "abc", "", "-1/2", rightAnswer(t, e))

	for _, want := range []string{
		"Não percebi. Escreve uma fração, por exemplo 3/4. Tenta outra vez.",
		"Escreve uma resposta. Tenta outra vez.",
		"Aqui não há números negativos. Tenta outra vez.",
		"Resultado: 1 de 1 certas.",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
}

func TestWrongAnswerShowsTheRightOne(t *testing.T) {
	e := exercises(t, exact.SkillFractionOfQuantity, 1, 1, 3)[0]
	right := rightAnswer(t, e)
	wrong := right + "1"

	out := play(t, []string{"-skill", exact.SkillFractionOfQuantity, "-n", "1", "-seed", "3"}, wrong)

	if want := "Ainda não. A resposta certa é " + right + "."; !strings.Contains(out, want) {
		t.Errorf("missing %q in:\n%s", want, out)
	}
	if !strings.Contains(out, "Resultado: 0 de 1 certas.") {
		t.Errorf("want 0 of 1 in:\n%s", out)
	}
}

func TestNotSimplified(t *testing.T) {
	e := exercises(t, exact.SkillEquivalentFractions, 2, 1, 1)[0]
	// Doubling both terms of the irreducible answer is equivalent but not
	// simplified.
	n, d := e.Answer.Num().Int64(), e.Answer.Denom().Int64()
	typed := strconv.FormatInt(2*n, 10) + "/" + strconv.FormatInt(2*d, 10)

	out := play(t, []string{"-level", "2", "-n", "1", "-seed", "1"}, typed)

	for _, want := range []string{
		"Quase! Está certo, mas falta simplificar. A fração irredutível é " + e.Answer.RatString() + ".",
		"Resultado: 0 de 1 certas. Mais 1 certas mas por simplificar.",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
}

func TestMisconceptionHint(t *testing.T) {
	// Try seeds until each misconception has been chosen once.
	seen := map[exact.Misconception]bool{}
	for seed := uint64(1); len(seen) < len(misconceptionHints) && seed < 50; seed++ {
		e := exercises(t, exact.SkillEquivalentFractions, 3, 1, seed)[0]
		for i, o := range e.Options {
			if o.Correct || seen[o.Misconception] {
				continue
			}
			args := []string{"-level", "3", "-n", "1", "-seed", strconv.FormatUint(seed, 10)}
			out := play(t, args, strconv.Itoa(i+1))

			if !strings.Contains(out, misconceptionHints[o.Misconception]) {
				t.Errorf("seed %d, chose %q: missing hint for %s in:\n%s", seed, o.Text, o.Misconception, out)
			}
			if want := "A resposta certa é a opção " + strconv.Itoa(correctOption(e)); !strings.Contains(out, want) {
				t.Errorf("seed %d: missing %q in:\n%s", seed, want, out)
			}
			seen[o.Misconception] = true
		}
	}
	if len(seen) != len(misconceptionHints) {
		t.Errorf("only exercised %v", seen)
	}
}

func TestEndOfInputStopsCleanly(t *testing.T) {
	e := exercises(t, exact.SkillEquivalentFractions, 1, 3, 9)[0]
	var out bytes.Buffer
	in := strings.NewReader(rightAnswer(t, e) + "\n") // answers one of three, then EOF

	if err := run([]string{"-n", "3", "-seed", "9"}, in, &out, &bytes.Buffer{}); err != nil {
		t.Fatalf("run: %v", err)
	}
	for _, want := range []string{"Sessão interrompida.", "Resultado: 1 de 1 certas."} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("missing %q in:\n%s", want, out.String())
		}
	}
}

func TestRandomSeedIsPrintedAndReplays(t *testing.T) {
	out := play(t, []string{"-n", "1"}, "1")

	_, rest, ok := strings.Cut(out, "(seed ")
	if !ok {
		t.Fatalf("no seed in:\n%s", out)
	}
	seed, _, _ := strings.Cut(rest, ")")

	replay := play(t, []string{"-n", "1", "-seed", seed}, "1")
	if replay != out {
		t.Errorf("replay with -seed %s differs:\n%s\nvs\n%s", seed, out, replay)
	}
}

func TestBadFlags(t *testing.T) {
	tests := []struct {
		args []string
		want string
	}{
		{[]string{"-skill", "mat.nao_existe"}, "available: mat.fracao_de_quantidade"},
		{[]string{"-level", "0"}, "has levels 1–3"},
		{[]string{"-level", "4"}, "has levels 1–3"},
		{[]string{"-n", "0"}, "want 1–50"},
		{[]string{"-n", "51"}, "want 1–50"},
		{[]string{"extra"}, "unexpected arguments"},
		{[]string{"-seed", "-1"}, "invalid value"},
	}
	for _, tt := range tests {
		err := run(tt.args, strings.NewReader(""), &bytes.Buffer{}, &bytes.Buffer{})
		if err == nil || !strings.Contains(err.Error(), tt.want) {
			t.Errorf("run(%v) error = %v, want it to contain %q", tt.args, err, tt.want)
		}
	}
}
