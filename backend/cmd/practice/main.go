// Command practice runs a short maths practice session in the terminal:
// one exercise at a time, immediate feedback in Portuguese, and a score at
// the end.
//
//	go run ./cmd/practice -skill mat.fracoes_equivalentes -level 2 -n 5 -seed 42
//
// Without -seed it picks a random one and prints it, so any session can be
// replayed exactly.
package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"math/rand/v2"
	"os"
	"strings"

	"github.com/WHAIT666/tpc-plus/backend/internal/engine/exact"
	"github.com/WHAIT666/tpc-plus/backend/internal/grading"
)

const maxExercises = 50

func main() {
	err := run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr)
	switch {
	case err == nil:
	case errors.Is(err, flag.ErrHelp):
		// -h: the flag package already printed the usage.
	default:
		fmt.Fprintln(os.Stderr, "practice:", err)
		os.Exit(1)
	}
}

type config struct {
	skill   string
	level   int
	n       int
	seed    uint64
	seedSet bool
}

func parseFlags(args []string, stderr io.Writer) (config, error) {
	var c config
	fs := flag.NewFlagSet("practice", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.StringVar(&c.skill, "skill", exact.SkillEquivalentFractions, "skill ID from the catalogue")
	fs.IntVar(&c.level, "level", 1, "difficulty level, from 1")
	fs.IntVar(&c.n, "n", 5, "number of exercises")
	fs.Uint64Var(&c.seed, "seed", 0, "replay the session with this seed (default: random)")

	if err := fs.Parse(args); err != nil {
		return config{}, err
	}
	if fs.NArg() > 0 {
		return config{}, fmt.Errorf("unexpected arguments: %s", strings.Join(fs.Args(), " "))
	}
	fs.Visit(func(f *flag.Flag) {
		if f.Name == "seed" {
			c.seedSet = true // so -seed 0 is a real seed, not "random"
		}
	})
	if c.n < 1 || c.n > maxExercises {
		return config{}, fmt.Errorf("-n %d: want 1–%d", c.n, maxExercises)
	}
	return c, nil
}

// run is main without the process: tests drive it with their own input and
// read its output.
func run(args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	c, err := parseFlags(args, stderr)
	if err != nil {
		return err
	}

	reg, err := exact.Default()
	if err != nil {
		return err
	}
	g, err := reg.Generator(c.skill)
	if err != nil {
		return fmt.Errorf("%w (available: %s)", err, strings.Join(reg.SkillIDs(), ", "))
	}
	if c.level < 1 || c.level > g.Levels() {
		return fmt.Errorf("-level %d: %s has levels 1–%d", c.level, c.skill, g.Levels())
	}
	if !c.seedSet {
		// Small enough to read aloud and type back.
		c.seed = rand.Uint64N(1_000_000) + 1
	}

	s := session{cfg: c, gen: g, in: bufio.NewScanner(stdin), out: stdout}
	return s.run()
}

type session struct {
	cfg config
	gen exact.Generator
	in  *bufio.Scanner
	out io.Writer

	answered, correct, notSimplified int
}

func (s *session) run() error {
	s.printf("Prática: %s, nível %d (seed %d)\n", s.cfg.skill, s.cfg.level, s.cfg.seed)

	for i := range s.cfg.n {
		// Exercise i uses seed+i, so the whole session follows from one seed.
		e, err := s.gen.Generate(s.cfg.seed+uint64(i), s.cfg.level)
		if err != nil {
			return err
		}

		s.printf("\nExercício %d de %d\n%s\n", i+1, s.cfg.n, e.Prompt)
		for j, o := range e.Options {
			s.printf("  %d) %s\n", j+1, o.Text)
		}

		typed, res, err := s.ask(e)
		if errors.Is(err, io.EOF) {
			s.printf("\nSessão interrompida.\n")
			break
		}
		if err != nil {
			return err
		}

		s.answered++
		switch res.Outcome {
		case grading.OutcomeCorrect:
			s.correct++
		case grading.OutcomeNotSimplified:
			s.notSimplified++
		}
		s.printf("%s\n", feedback(e, typed, res))
	}

	s.printScore()
	return nil
}

// ask reads answers until one can be graded. Input the parser rejects
// does not count as an attempt: the child gets a hint and tries again.
func (s *session) ask(e exact.Exercise) (string, grading.Result, error) {
	key := e.Key()
	for {
		s.printf("> ")
		if !s.in.Scan() {
			if err := s.in.Err(); err != nil {
				return "", grading.Result{}, fmt.Errorf("read answer: %w", err)
			}
			return "", grading.Result{}, io.EOF
		}
		typed := s.in.Text()

		res, err := grading.Grade(typed, key)
		if err != nil {
			return "", grading.Result{}, fmt.Errorf("exercise %s seed %d level %d: %w", e.SkillID, e.Seed, e.Level, err)
		}
		if res.Outcome != grading.OutcomeInvalidInput {
			return typed, res, nil
		}
		s.printf("%s Tenta outra vez.\n", invalidHint(e, res.Err))
	}
}

func (s *session) printScore() {
	s.printf("\nResultado: %d de %d certas.", s.correct, s.answered)
	if s.notSimplified > 0 {
		s.printf(" Mais %d certas mas por simplificar.", s.notSimplified)
	}
	s.printf("\nPara repetir: go run ./cmd/practice -skill %s -level %d -n %d -seed %d\n",
		s.cfg.skill, s.cfg.level, s.cfg.n, s.cfg.seed)
}

// printf ignores write errors: a terminal that cannot be written to has no
// one reading it either.
func (s *session) printf(format string, args ...any) {
	_, _ = fmt.Fprintf(s.out, format, args...)
}
