package builtin

import "github.com/orchestra-mcp/themes/src/types"

// LightTheme returns the built-in orchestra-light theme.
func LightTheme() *types.ThemeDef {
	return &types.ThemeDef{
		ID:          "orchestra-light",
		Name:        "Orchestra Light",
		Description: "Default light theme for Orchestra",
		Author:      "Orchestra Team",
		Type:        "light",
		Colors:      lightColors(),
	}
}

// DarkTheme returns the built-in orchestra-dark theme.
func DarkTheme() *types.ThemeDef {
	return &types.ThemeDef{
		ID:          "orchestra-dark",
		Name:        "Orchestra Dark",
		Description: "Default dark theme for Orchestra",
		Author:      "Orchestra Team",
		Type:        "dark",
		Colors:      darkColors(),
	}
}

func lightColors() map[string]string {
	return map[string]string{
		"background":           "#FFFFFF",
		"foreground":           "#1E1E1E",
		"primary":             "#2563EB",
		"secondary":           "#64748B",
		"accent":              "#8B5CF6",
		"border":              "#E2E8F0",
		"sidebar.background":  "#F8FAFC",
		"sidebar.foreground":  "#334155",
		"editor.background":   "#FFFFFF",
		"editor.foreground":   "#1E1E1E",
		"titlebar.background": "#F1F5F9",
		"titlebar.foreground": "#0F172A",
		"statusbar.background": "#2563EB",
		"statusbar.foreground": "#FFFFFF",
		"input.background":    "#FFFFFF",
		"input.foreground":    "#1E1E1E",
		"input.border":        "#CBD5E1",
		"button.background":   "#2563EB",
		"button.foreground":   "#FFFFFF",
		"error":               "#DC2626",
		"warning":             "#D97706",
		"success":             "#16A34A",
		"info":                "#2563EB",
	}
}

func darkColors() map[string]string {
	return map[string]string{
		"background":           "#1E1E2E",
		"foreground":           "#CDD6F4",
		"primary":             "#89B4FA",
		"secondary":           "#A6ADC8",
		"accent":              "#CBA6F7",
		"border":              "#313244",
		"sidebar.background":  "#181825",
		"sidebar.foreground":  "#BAC2DE",
		"editor.background":   "#1E1E2E",
		"editor.foreground":   "#CDD6F4",
		"titlebar.background": "#11111B",
		"titlebar.foreground": "#CDD6F4",
		"statusbar.background": "#89B4FA",
		"statusbar.foreground": "#1E1E2E",
		"input.background":    "#313244",
		"input.foreground":    "#CDD6F4",
		"input.border":        "#45475A",
		"button.background":   "#89B4FA",
		"button.foreground":   "#1E1E2E",
		"error":               "#F38BA8",
		"warning":             "#FAB387",
		"success":             "#A6E3A1",
		"info":                "#89B4FA",
	}
}

// BuiltinThemes returns all built-in themes.
func BuiltinThemes() []*types.ThemeDef {
	return []*types.ThemeDef{LightTheme(), DarkTheme()}
}
