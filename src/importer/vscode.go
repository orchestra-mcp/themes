package importer

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/orchestra-mcp/themes/src/types"
)

// vscodeColorMap maps VS Code color keys to Orchestra color keys.
var vscodeColorMap = map[string]string{
	"editor.background":                "bg-primary",
	"editor.foreground":                "text-primary",
	"sideBar.background":               "bg-secondary",
	"activityBar.background":           "bg-tertiary",
	"statusBar.background":             "bg-accent",
	"titleBar.activeBackground":        "bg-header",
	"focusBorder":                      "border-focus",
	"list.activeSelectionBackground":   "bg-selection",
}

// vscodeThemeFile represents the top-level VS Code theme JSON structure.
type vscodeThemeFile struct {
	Name        string                  `json:"name"`
	Type        string                  `json:"type"`
	Include     string                  `json:"include,omitempty"`
	Colors      map[string]string       `json:"colors"`
	TokenColors []vscodeTokenColorEntry `json:"tokenColors"`
}

// vscodeTokenColorEntry represents a single token color rule.
type vscodeTokenColorEntry struct {
	Name     string              `json:"name"`
	Scope    vscodeScope         `json:"scope"`
	Settings vscodeTokenSettings `json:"settings"`
}

// vscodeScope handles VS Code scope which can be string or []string.
type vscodeScope []string

// UnmarshalJSON handles both string and string array scope values.
func (s *vscodeScope) UnmarshalJSON(data []byte) error {
	var single string
	if err := json.Unmarshal(data, &single); err == nil {
		*s = splitScope(single)
		return nil
	}
	var arr []string
	if err := json.Unmarshal(data, &arr); err != nil {
		return fmt.Errorf("scope must be string or string array: %w", err)
	}
	*s = arr
	return nil
}

// splitScope splits a comma-separated scope string.
func splitScope(raw string) []string {
	parts := strings.Split(raw, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// vscodeTokenSettings holds foreground and fontStyle for a token rule.
type vscodeTokenSettings struct {
	Foreground string `json:"foreground"`
	FontStyle  string `json:"fontStyle"`
}

// ImportVSCode parses a VS Code JSON theme file into a ThemeDef.
func ImportVSCode(data []byte) (*types.ThemeDef, error) {
	var vsTheme vscodeThemeFile
	if err := json.Unmarshal(data, &vsTheme); err != nil {
		return nil, fmt.Errorf("invalid VS Code theme JSON: %w", err)
	}

	theme := &types.ThemeDef{
		Name:   vsTheme.Name,
		Type:   normalizeThemeType(vsTheme.Type),
		Source: "vscode",
		Colors: make(map[string]string),
	}

	if vsTheme.Include != "" {
		theme.Colors["_include"] = vsTheme.Include
	}

	mapVSCodeColors(vsTheme.Colors, theme.Colors)
	theme.TokenColors = mapVSCodeTokenColors(vsTheme.TokenColors)

	if theme.Name == "" {
		theme.Name = "Imported VS Code Theme"
	}
	theme.ID = slugify(theme.Name)

	return theme, nil
}

// mapVSCodeColors maps VS Code color keys to Orchestra keys
// and stores unmapped keys under the "raw." prefix.
func mapVSCodeColors(src map[string]string, dst map[string]string) {
	for key, value := range src {
		if mapped, ok := vscodeColorMap[key]; ok {
			dst[mapped] = value
		}
		dst["raw."+key] = value
	}
}

// mapVSCodeTokenColors converts VS Code token colors to Orchestra format.
func mapVSCodeTokenColors(entries []vscodeTokenColorEntry) []types.TokenColor {
	if len(entries) == 0 {
		return nil
	}
	result := make([]types.TokenColor, 0, len(entries))
	for _, entry := range entries {
		settings := make(map[string]string)
		if entry.Settings.Foreground != "" {
			settings["foreground"] = entry.Settings.Foreground
		}
		if entry.Settings.FontStyle != "" {
			settings["fontStyle"] = entry.Settings.FontStyle
		}
		if len(settings) == 0 && len(entry.Scope) == 0 {
			continue
		}
		tc := types.TokenColor{
			Name:     entry.Name,
			Scope:    entry.Scope,
			Settings: settings,
		}
		result = append(result, tc)
	}
	return result
}

// normalizeThemeType converts VS Code theme type to Orchestra type.
func normalizeThemeType(vsType string) string {
	switch strings.ToLower(vsType) {
	case "dark", "vs-dark":
		return "dark"
	case "light", "vs":
		return "light"
	case "hc-black", "hc-light":
		return "high-contrast"
	default:
		return "dark"
	}
}

// slugify creates a URL-safe ID from a theme name.
func slugify(name string) string {
	s := strings.ToLower(name)
	s = strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			return r
		}
		if r == ' ' || r == '-' || r == '_' {
			return '-'
		}
		return -1
	}, s)
	for strings.Contains(s, "--") {
		s = strings.ReplaceAll(s, "--", "-")
	}
	return strings.Trim(s, "-")
}
