package tests

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/orchestra-mcp/themes/src/builtin"
	"github.com/orchestra-mcp/themes/src/service"
	"github.com/orchestra-mcp/themes/src/types"
	"github.com/rs/zerolog"
)

func newTestService(t *testing.T) *service.ThemesService {
	t.Helper()
	dir := t.TempDir()
	logger := zerolog.Nop()
	return service.New(dir, "orchestra-dark", logger)
}

func TestBuiltinThemesLoaded(t *testing.T) {
	svc := newTestService(t)
	themes := svc.GetAvailableThemes()
	if len(themes) < 2 {
		t.Fatalf("expected at least 2 built-in themes, got %d", len(themes))
	}

	light, err := svc.GetTheme("orchestra-light")
	if err != nil {
		t.Fatalf("expected orchestra-light: %v", err)
	}
	if light.Type != "light" {
		t.Errorf("expected type 'light', got %q", light.Type)
	}

	dark, err := svc.GetTheme("orchestra-dark")
	if err != nil {
		t.Fatalf("expected orchestra-dark: %v", err)
	}
	if dark.Type != "dark" {
		t.Errorf("expected type 'dark', got %q", dark.Type)
	}
}

func TestDefaultActiveTheme(t *testing.T) {
	svc := newTestService(t)
	active := svc.GetActiveTheme()
	if active == nil {
		t.Fatal("expected active theme, got nil")
	}
	if active.ID != "orchestra-dark" {
		t.Errorf("expected default 'orchestra-dark', got %q", active.ID)
	}
}

func TestSetActiveTheme(t *testing.T) {
	svc := newTestService(t)

	err := svc.SetActiveTheme("orchestra-light")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	active := svc.GetActiveTheme()
	if active.ID != "orchestra-light" {
		t.Errorf("expected 'orchestra-light', got %q", active.ID)
	}
}

func TestSetActiveThemeNotFound(t *testing.T) {
	svc := newTestService(t)
	err := svc.SetActiveTheme("nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent theme")
	}
}

func TestRegisterTheme(t *testing.T) {
	svc := newTestService(t)
	custom := &types.ThemeDef{
		ID:   "custom-theme",
		Name: "Custom Theme",
		Type: "dark",
		Colors: map[string]string{
			"background": "#000000",
			"foreground": "#FFFFFF",
		},
	}
	svc.RegisterTheme(custom)

	got, err := svc.GetTheme("custom-theme")
	if err != nil {
		t.Fatalf("expected custom theme: %v", err)
	}
	if got.Name != "Custom Theme" {
		t.Errorf("expected name 'Custom Theme', got %q", got.Name)
	}
}

func TestChangeListener(t *testing.T) {
	svc := newTestService(t)

	var oldID, newID string
	svc.OnDidChangeTheme(func(old, new *types.ThemeDef) {
		if old != nil {
			oldID = old.ID
		}
		if new != nil {
			newID = new.ID
		}
	})

	err := svc.SetActiveTheme("orchestra-light")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if oldID != "orchestra-dark" {
		t.Errorf("expected old theme 'orchestra-dark', got %q", oldID)
	}
	if newID != "orchestra-light" {
		t.Errorf("expected new theme 'orchestra-light', got %q", newID)
	}
}

func TestExportTheme(t *testing.T) {
	svc := newTestService(t)

	data, err := svc.ExportTheme("orchestra-dark")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var theme types.ThemeDef
	if err := json.Unmarshal(data, &theme); err != nil {
		t.Fatalf("invalid JSON export: %v", err)
	}
	if theme.ID != "orchestra-dark" {
		t.Errorf("expected ID 'orchestra-dark', got %q", theme.ID)
	}
}

func TestExportThemeNotFound(t *testing.T) {
	svc := newTestService(t)
	_, err := svc.ExportTheme("nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent theme export")
	}
}

func TestImportTheme(t *testing.T) {
	svc := newTestService(t)

	theme := types.ThemeDef{
		ID:     "imported-theme",
		Name:   "Imported Theme",
		Author: "Test Author",
		Type:   "light",
		Colors: map[string]string{
			"background": "#FAFAFA",
			"foreground": "#333333",
		},
	}
	data, _ := json.Marshal(theme)

	imported, err := svc.ImportTheme(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if imported.ID != "imported-theme" {
		t.Errorf("expected ID 'imported-theme', got %q", imported.ID)
	}

	got, err := svc.GetTheme("imported-theme")
	if err != nil {
		t.Fatalf("expected imported theme: %v", err)
	}
	if got.Author != "Test Author" {
		t.Errorf("expected author 'Test Author', got %q", got.Author)
	}
}

func TestImportThemeInvalidJSON(t *testing.T) {
	svc := newTestService(t)
	_, err := svc.ImportTheme([]byte("not json"))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestImportThemeMissingID(t *testing.T) {
	svc := newTestService(t)
	_, err := svc.ImportTheme([]byte(`{"name":"No ID"}`))
	if err == nil {
		t.Fatal("expected error for missing theme ID")
	}
}

func TestPreferencePersistence(t *testing.T) {
	dir := t.TempDir()
	logger := zerolog.Nop()

	svc1 := service.New(dir, "orchestra-dark", logger)
	if err := svc1.SetActiveTheme("orchestra-light"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify preference file exists.
	prefPath := filepath.Join(dir, "theme-preference.json")
	if _, err := os.Stat(prefPath); os.IsNotExist(err) {
		t.Fatal("expected preference file to exist")
	}

	// New instance should load persisted preference.
	svc2 := service.New(dir, "orchestra-dark", logger)
	active := svc2.GetActiveTheme()
	if active.ID != "orchestra-light" {
		t.Errorf("expected persisted 'orchestra-light', got %q", active.ID)
	}
}

func TestBuiltinThemeColorKeys(t *testing.T) {
	expectedKeys := []string{
		"background", "foreground", "primary", "secondary", "accent",
		"border", "sidebar.background", "sidebar.foreground",
		"editor.background", "editor.foreground",
		"titlebar.background", "titlebar.foreground",
		"statusbar.background", "statusbar.foreground",
		"input.background", "input.foreground", "input.border",
		"button.background", "button.foreground",
		"error", "warning", "success", "info",
	}

	for _, theme := range builtin.BuiltinThemes() {
		for _, key := range expectedKeys {
			if _, ok := theme.Colors[key]; !ok {
				t.Errorf("theme %q missing color key %q", theme.ID, key)
			}
		}
	}
}
