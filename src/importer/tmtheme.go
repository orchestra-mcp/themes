package importer

import (
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/orchestra-mcp/themes/src/types"
)

const (
	themeDark  = "dark"
	themeLight = "light"
)

// tmTheme global color mappings to Orchestra keys.
var tmGlobalColorMap = map[string]string{
	"background":    "bg-primary",
	"foreground":    "text-primary",
	"caret":         "caret",
	"selection":     "bg-selection",
	"lineHighlight": "bg-line-highlight",
}

// ImportTmTheme parses a .tmTheme plist XML file into a ThemeDef.
func ImportTmTheme(data []byte) (*types.ThemeDef, error) {
	var root plistRoot
	if err := xml.Unmarshal(data, &root); err != nil {
		return nil, fmt.Errorf("invalid tmTheme XML: %w", err)
	}

	theme := &types.ThemeDef{
		Source: "tmtheme",
		Colors: make(map[string]string),
	}

	theme.Name = dictGet(&root.Dict, "name")
	if theme.Name == "" {
		theme.Name = "Imported TextMate Theme"
	}
	theme.ID = slugify(theme.Name)
	theme.Author = dictGet(&root.Dict, "author")

	settings := dictGetArray(&root.Dict, "settings")
	if len(settings) == 0 {
		return nil, fmt.Errorf("tmTheme has no settings array")
	}

	extractGlobalColors(settings, theme)
	theme.TokenColors = extractTokenColors(settings)
	theme.Type = detectThemeType(theme.Colors)

	return theme, nil
}

// extractGlobalColors reads the first settings entry (global colors).
func extractGlobalColors(settings []plistDict, theme *types.ThemeDef) {
	if len(settings) == 0 {
		return
	}

	global := dictGetDict(&settings[0], "settings")
	if global == nil {
		return
	}

	for _, item := range global.Items {
		value := item.Value.String
		if value == "" {
			continue
		}
		if mapped, ok := tmGlobalColorMap[item.Key]; ok {
			theme.Colors[mapped] = value
		}
		theme.Colors["raw."+item.Key] = value
	}
}

// extractTokenColors reads token color rules from settings[1:].
func extractTokenColors(settings []plistDict) []types.TokenColor {
	if len(settings) <= 1 {
		return nil
	}

	result := make([]types.TokenColor, 0, len(settings)-1)
	for _, entry := range settings[1:] {
		name := dictGet(&entry, "name")
		scopeStr := dictGet(&entry, "scope")
		settingsDict := dictGetDict(&entry, "settings")

		if settingsDict == nil {
			continue
		}

		settingsMap := make(map[string]string)
		for _, item := range settingsDict.Items {
			if item.Value.String != "" {
				settingsMap[item.Key] = item.Value.String
			}
		}

		var scopes []string
		if scopeStr != "" {
			scopes = splitScope(scopeStr)
		}

		tc := types.TokenColor{
			Name:     name,
			Scope:    scopes,
			Settings: settingsMap,
		}
		result = append(result, tc)
	}
	return result
}

// detectThemeType guesses light vs dark from the background color.
func detectThemeType(colors map[string]string) string {
	bg := colors["bg-primary"]
	if bg == "" {
		bg = colors["raw.background"]
	}
	if bg == "" {
		return themeDark
	}
	if isLightColor(bg) {
		return themeLight
	}
	return themeDark
}

// isLightColor returns true if a hex color is considered light.
func isLightColor(hex string) bool {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) < 6 {
		return false
	}
	r := hexVal(hex[0])*16 + hexVal(hex[1])
	g := hexVal(hex[2])*16 + hexVal(hex[3])
	b := hexVal(hex[4])*16 + hexVal(hex[5])
	luminance := (r*299 + g*587 + b*114) / 1000
	return luminance > 128
}

func hexVal(c byte) int {
	switch {
	case c >= '0' && c <= '9':
		return int(c - '0')
	case c >= 'a' && c <= 'f':
		return int(c-'a') + 10
	case c >= 'A' && c <= 'F':
		return int(c-'A') + 10
	default:
		return 0
	}
}
