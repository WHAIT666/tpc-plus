package exact

import (
	"errors"
	"fmt"
	"slices"
)

var (
	// ErrUnknownSkill is returned for a skill ID with no generator.
	ErrUnknownSkill = errors.New("unknown skill")

	// ErrDuplicateSkill is returned when two generators share a skill ID.
	ErrDuplicateSkill = errors.New("duplicate skill")
)

// Registry finds generators by skill ID.
type Registry struct {
	gens map[string]Generator
}

// NewRegistry returns a registry of gens. Skill IDs must be unique.
func NewRegistry(gens ...Generator) (*Registry, error) {
	r := &Registry{gens: make(map[string]Generator, len(gens))}
	for _, g := range gens {
		id := g.SkillID()
		if _, ok := r.gens[id]; ok {
			return nil, fmt.Errorf("new registry: %q: %w", id, ErrDuplicateSkill)
		}
		r.gens[id] = g
	}
	return r, nil
}

// Default returns a registry with every generator in this package.
func Default() (*Registry, error) {
	return NewRegistry(
		equivalentFractions{},
	)
}

// Generator returns the generator for a skill ID.
func (r *Registry) Generator(skillID string) (Generator, error) {
	g, ok := r.gens[skillID]
	if !ok {
		return nil, fmt.Errorf("generator %q: %w", skillID, ErrUnknownSkill)
	}
	return g, nil
}

// SkillIDs returns every registered skill ID, sorted.
func (r *Registry) SkillIDs() []string {
	ids := make([]string, 0, len(r.gens))
	for id := range r.gens {
		ids = append(ids, id)
	}
	slices.Sort(ids)
	return ids
}
