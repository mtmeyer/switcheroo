package main

import "testing"

func TestRunConfiguredCommand_ReplacesPath(t *testing.T) {
	originalRunner := execCommandRunner
	defer func() { execCommandRunner = originalRunner }()

	var gotName string
	var gotArgs []string
	execCommandRunner = func(name string, args []string) error {
		gotName = name
		gotArgs = append([]string(nil), args...)
		return nil
	}

	if err := runConfiguredCommand("echo {{path}}", "/tmp/example"); err != nil {
		t.Fatalf("runConfiguredCommand error: %v", err)
	}

	if gotName != "echo" {
		t.Fatalf("expected command name 'echo', got %q", gotName)
	}
	if len(gotArgs) != 1 || gotArgs[0] != "/tmp/example" {
		t.Fatalf("unexpected args: %v", gotArgs)
	}
}

func TestRunConfiguredCommand_InvalidTemplate(t *testing.T) {
	if err := runConfiguredCommand("   ", "/tmp/example"); err == nil {
		t.Fatal("expected error for empty command template")
	}
}
