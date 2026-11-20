package exact

import (
	"errors"
	"reflect"
	"slices"
	"testing"
)

func TestDefaultRegistry(t *testing.T) {
	reg, err := Default()
	if err != nil {
		t.Fatalf("Default(): %v", err)
	}

	want := []string{SkillFractionOfQuantity, SkillEquivalentFractions, SkillAddSubtractFractions}
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
		if !reflect.DeepEqual(a, b) {
			t.Errorf("level %d: seed 42 gave %+v, then %+v", level, a, b)
		}
	}
}
