package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempConfig(t *testing.T, contents string) string {
	t.Helper()
	file, err := os.CreateTemp(t.TempDir(), "switcheroo-config-*.json")
	if err != nil {
		t.Fatalf("create temp config: %v", err)
	}
	if _, err := file.WriteString(contents); err != nil {
		file.Close()
		t.Fatalf("write temp config: %v", err)
	}
	if err := file.Close(); err != nil {
		t.Fatalf("close temp config: %v", err)
	}
	return file.Name()
}

func TestLoadConfig_DefaultOutput(t *testing.T) {
	path := writeTempConfig(t, `{
	  "directory": "`+filepath.ToSlash(t.TempDir())+`"
	}`)

	cfg, err := LoadFromPath(path)
	if err != nil {
		t.Fatalf("LoadFromPath error: %v", err)
	}
	if cfg.Output.Type != "path" {
		t.Fatalf("expected default output.type=path, got %q", cfg.Output.Type)
	}
}

func TestLoadConfig_CommandRequiresValue(t *testing.T) {
	path := writeTempConfig(t, `{
	  "directory": "`+filepath.ToSlash(t.TempDir())+`",
	  "output": {
	    "type": "command"
	  }
	}`)

	if _, err := LoadFromPath(path); err == nil {
		t.Fatal("expected error when command output missing value")
	}
}

func TestLoadConfig_CommandValid(t *testing.T) {
	path := writeTempConfig(t, `{
	  "directory": "`+filepath.ToSlash(t.TempDir())+`",
	  "output": {
	    "type": "command",
	    "value": "echo {{path}}"
	  }
	}`)

	cfg, err := LoadFromPath(path)
	if err != nil {
		t.Fatalf("LoadFromPath error: %v", err)
	}
	if cfg.Output.Type != "command" {
		t.Fatalf("expected output.type=command, got %q", cfg.Output.Type)
	}
	if cfg.Output.Value != "echo {{path}}" {
		t.Fatalf("unexpected output.value: %q", cfg.Output.Value)
	}
}

func TestLoadConfig_UnsupportedOutputType(t *testing.T) {
	path := writeTempConfig(t, `{
	  "directory": "`+filepath.ToSlash(t.TempDir())+`",
	  "output": {
	    "type": "unknown"
	  }
	}`)

	if _, err := LoadFromPath(path); err == nil {
		t.Fatal("expected error for unsupported output.type")
	}
}
