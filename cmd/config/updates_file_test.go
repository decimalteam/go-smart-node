package config

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestSaveShorterUpgradeJournal(t *testing.T) {
	file := filepath.Join(t.TempDir(), "updates.json")
	plan := NewUpdatesInfo(file)
	plan.AddExecutedPlan(strings.Repeat("old-plan", 40), 99999999)
	if err := plan.Save(); err != nil {
		t.Fatal(err)
	}
	plan.AllBlocks = map[string]int64{}
	if err := plan.Save(); err != nil {
		t.Fatal(err)
	}
	loaded := NewUpdatesInfo(file)
	if err := loaded.Load(); err != nil {
		t.Fatalf("shorter rewrite left invalid JSON: %v", err)
	}
	if len(loaded.AllBlocks) != 0 {
		t.Fatal("old journal contents survived replacement")
	}
}
