package exact

import (
	"testing"

	"github.com/WHAIT666/tpc-plus/backend/internal/grading"
)

// assertSelfGradesCorrect generates one exercise, types its canonical answer
// the way the child would and checks that the grader says correct.
func assertSelfGradesCorrect(t *testing.T, g Generator, seed uint64, level int) {
	t.Helper()

	e, err := g.Generate(seed, level)
	if err != nil {
		t.Fatalf("%s Generate(%d, %d): %v", g.SkillID(), seed, level, err)
	}

	typed, err := grading.Format(e.Answer, e.Kind)
	if err != nil {
		t.Fatalf("%s seed %d level %d: format answer: %v", g.SkillID(), seed, level, err)
	}
	res, err := grading.Grade(typed, e.Key())
	if err != nil {
		t.Fatalf("%s seed %d level %d: invalid key: %v", g.SkillID(), seed, level, err)
	}
	if res.Outcome != grading.OutcomeCorrect {
		t.Fatalf("%s seed %d level %d: own answer %q graded %v", g.SkillID(), seed, level, typed, res.Outcome)
	}
}
