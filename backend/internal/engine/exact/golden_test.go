package exact

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

var update = flag.Bool("update", false, "rewrite the golden files in testdata/golden")

// goldenDir is backend/testdata/golden, seen from this package.
const goldenDir = "../../../testdata/golden"

// goldenSeeds are pinned for every skill and level.
const goldenSeeds = 20

// TestGolden pins the exercises for seeds 1–20 of every skill and level.
// If it fails after an intentional change, bump the generator's version and
// rewrite the files:
//
//	go test ./internal/engine/exact -run TestGolden -update
func TestGolden(t *testing.T) {
	reg, err := Default()
	if err != nil {
		t.Fatal(err)
	}

	for _, id := range reg.SkillIDs() {
		g, err := reg.Generator(id)
		if err != nil {
			t.Fatal(err)
		}
		for level := 1; level <= g.Levels(); level++ {
			t.Run(fmt.Sprintf("%s/level%d", id, level), func(t *testing.T) {
				got := goldenJSON(t, g, level)
				path := filepath.Join(goldenDir, id, fmt.Sprintf("level%d.json", level))

				if *update {
					if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
						t.Fatal(err)
					}
					if err := os.WriteFile(path, got, 0o644); err != nil {
						t.Fatal(err)
					}
					return
				}

				want, err := os.ReadFile(path)
				if err != nil {
					t.Fatalf("read golden file (run with -update to create it): %v", err)
				}
				if !bytes.Equal(got, want) {
					t.Errorf("output for %s level %d changed; if intended, bump the version and run with -update\n got: %s\nwant: %s",
						id, level, got, want)
				}
			})
		}
	}
}

func goldenJSON(t *testing.T, g Generator, level int) []byte {
	t.Helper()

	exercises := make([]Exercise, 0, goldenSeeds)
	for seed := uint64(1); seed <= goldenSeeds; seed++ {
		e, err := g.Generate(seed, level)
		if err != nil {
			t.Fatalf("Generate(%d, %d): %v", seed, level, err)
		}
		exercises = append(exercises, e)
	}

	b, err := json.MarshalIndent(exercises, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	return append(b, '\n')
}
