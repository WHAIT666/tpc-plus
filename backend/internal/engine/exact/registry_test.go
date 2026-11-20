package exact

import (
	"errors"
	"slices"
	"testing"
)

func TestDefaultRegistry(t *testing.T) {
	reg, err := Default()
	if err != nil {
		t.Fatalf("Default(): %v", err)
	}

	want := []string{SkillEquivalentFractions}
	if got := reg.SkillIDs(); !slices.Equal(got, want) {
		t.Errorf("SkillIDs() = %v, want %v", got, want)
	}

	g, err := reg.Generator(SkillEquivalentFractions)
	if err != nil {
		t.Fatalf("Generator(%q): %v", SkillEquivalentFractions, err)
	}
	if g.SkillID() != SkillEquivalentFractions {
		t.Errorf("Generator(%q).SkillID() = %q", SkillEquivalentFractions, g.SkillID())
	}
}

func TestRegistryUnknownSkill(t *testing.T) {
	reg, err := Default()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := reg.Generator("mat.nao_existe"); !errors.Is(err, ErrUnknownSkill) {
		t.Errorf("error = %v, want %v", err, ErrUnknownSkill)
	}
}

func TestRegistryDuplicateSkill(t *testing.T) {
	_, err := NewRegistry(equivalentFractions{}, equivalentFractions{})
	if !errors.Is(err, ErrDuplicateSkill) {
		t.Errorf("error = %v, want %v", err, ErrDuplicateSkill)
	}
}

// Every exercise from every generator must be gradable, and its own answer
// must be graded correct. The property tests in step 6 check the values
// against an independent oracle; this only checks the wiring.
func TestEveryAnswerGradesCorrect(t *testing.T) {
	reg, err := Default()
	if err != nil {
		t.Fatal(err)
	}

	for _, id := range reg.SkillIDs() {
		g, _ := reg.Generator(id)
		for level := 1; level <= g.Levels(); level++ {
			for seed := uint64(1); seed <= 200; seed++ {
				assertSelfGradesCorrect(t, g, seed, level)
			}
		}
	}
}

func TestGenerateUnknownLevel(t *testing.T) {
	g := equivalentFractions{}
	for _, level := range []int{0, -1, g.Levels() + 1} {
		if _, err := g.Generate(1, level); !errors.Is(err, ErrUnknownLevel) {
			t.Errorf("Generate(1, %d) error = %v, want %v", level, err, ErrUnknownLevel)
		}
	}
}

func TestGenerateIsDeterministic(t *testing.T) {
	g := equivalentFractions{}
	for level := 1; level <= g.Levels(); level++ {
		a, _ := g.Generate(42, level)
		b, _ := g.Generate(42, level)
		if a.Prompt != b.Prompt || a.Answer.Cmp(b.Answer) != 0 {
			t.Errorf("level %d: seed 42 gave %q = %v, then %q = %v", level, a.Prompt, a.Answer, b.Prompt, b.Answer)
		}
	}
}
