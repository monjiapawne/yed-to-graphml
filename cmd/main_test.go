package main

import (
	"path/filepath"
	"testing"
)

func TestLoadNodesFromPath_File(t *testing.T) {
	file := filepath.Join("../test", "1_go.svg")
	m, err := loadNodesFromPath(file, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m) != 1 {
		t.Errorf("expected 1 group, got %d", len(m))
	}
}

func TestLoadNodesFromPath_Dir(t *testing.T) {
	// assumes test/ contains at least one svg
	m, err := loadNodesFromPath("../test", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	found := false
	for _, nodes := range m {
		if len(nodes) > 0 {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected at least one node in directory")
	}
}

func TestCfgValidate(t *testing.T) {
	cfg := Cfg{inPath: "../test", outPath: "../out"}
	if err := cfg.validate(); err != nil {
		t.Errorf("validate failed: %v", err)
	}
	cfg2 := Cfg{inPath: "", outPath: "../out"}
	if err := cfg2.validate(); err == nil {
		t.Error("expected error for empty inPath")
	}
}
