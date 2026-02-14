package tests

import (
	"testing"

	"github.com/orchestra-mcp/themes/src/importer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Format Detection ---

func TestDetectFormatVSCodeJSON(t *testing.T) {
	data := []byte(`{
		"name": "Monokai",
		"type": "dark",
		"colors": {
			"editor.background": "#272822"
		},
		"tokenColors": []
	}`)
	assert.Equal(t, importer.FormatVSCodeJSON, importer.DetectFormat(data))
}

func TestDetectFormatTmThemeXML(t *testing.T) {
	data := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0"><dict></dict></plist>`)
	assert.Equal(t, importer.FormatTmTheme, importer.DetectFormat(data))
}

func TestDetectFormatTmThemePlistPrefix(t *testing.T) {
	data := []byte(`<plist version="1.0"><dict></dict></plist>`)
	assert.Equal(t, importer.FormatTmTheme, importer.DetectFormat(data))
}

func TestDetectFormatOrchestraJSON(t *testing.T) {
	data := []byte(`{
		"id": "my-theme",
		"name": "My Theme",
		"type": "dark",
		"colors": {"background": "#000"}
	}`)
	assert.Equal(t, importer.FormatOrchestraJSON, importer.DetectFormat(data))
}

func TestDetectFormatFallback(t *testing.T) {
	data := []byte(`not valid anything`)
	assert.Equal(t, importer.FormatOrchestraJSON, importer.DetectFormat(data))
}

// --- VS Code JSON Import ---

var sampleVSCodeTheme = []byte(`{
	"name": "One Dark Pro",
	"type": "dark",
	"colors": {
		"editor.background": "#282c34",
		"editor.foreground": "#abb2bf",
		"sideBar.background": "#21252b",
		"activityBar.background": "#2c313a",
		"statusBar.background": "#21252b",
		"titleBar.activeBackground": "#282c34",
		"focusBorder": "#528bff",
		"list.activeSelectionBackground": "#2c313a",
		"tab.activeBackground": "#282c34"
	},
	"tokenColors": [
		{
			"name": "Comments",
			"scope": "comment",
			"settings": {
				"foreground": "#5c6370",
				"fontStyle": "italic"
			}
		},
		{
			"name": "Strings",
			"scope": ["string.quoted.double", "string.quoted.single"],
			"settings": {
				"foreground": "#98c379"
			}
		}
	]
}`)

func TestImportVSCodeBasicFields(t *testing.T) {
	theme, err := importer.ImportVSCode(sampleVSCodeTheme)
	require.NoError(t, err)

	assert.Equal(t, "one-dark-pro", theme.ID)
	assert.Equal(t, "One Dark Pro", theme.Name)
	assert.Equal(t, "dark", theme.Type)
	assert.Equal(t, "vscode", theme.Source)
}

func TestImportVSCodeColorMapping(t *testing.T) {
	theme, err := importer.ImportVSCode(sampleVSCodeTheme)
	require.NoError(t, err)

	assert.Equal(t, "#282c34", theme.Colors["bg-primary"])
	assert.Equal(t, "#abb2bf", theme.Colors["text-primary"])
	assert.Equal(t, "#21252b", theme.Colors["bg-secondary"])
	assert.Equal(t, "#2c313a", theme.Colors["bg-tertiary"])
	assert.Equal(t, "#21252b", theme.Colors["bg-accent"])
	assert.Equal(t, "#282c34", theme.Colors["bg-header"])
	assert.Equal(t, "#528bff", theme.Colors["border-focus"])
	assert.Equal(t, "#2c313a", theme.Colors["bg-selection"])
}

func TestImportVSCodeRawColors(t *testing.T) {
	theme, err := importer.ImportVSCode(sampleVSCodeTheme)
	require.NoError(t, err)

	assert.Equal(t, "#282c34", theme.Colors["raw.editor.background"])
	assert.Equal(t, "#282c34", theme.Colors["raw.tab.activeBackground"])
}

func TestImportVSCodeTokenColors(t *testing.T) {
	theme, err := importer.ImportVSCode(sampleVSCodeTheme)
	require.NoError(t, err)

	require.Len(t, theme.TokenColors, 2)

	assert.Equal(t, "Comments", theme.TokenColors[0].Name)
	assert.Equal(t, []string{"comment"}, theme.TokenColors[0].Scope)
	assert.Equal(t, "#5c6370", theme.TokenColors[0].Settings["foreground"])
	assert.Equal(t, "italic", theme.TokenColors[0].Settings["fontStyle"])

	assert.Equal(t, "Strings", theme.TokenColors[1].Name)
	assert.Contains(t, theme.TokenColors[1].Scope, "string.quoted.double")
	assert.Contains(t, theme.TokenColors[1].Scope, "string.quoted.single")
	assert.Equal(t, "#98c379", theme.TokenColors[1].Settings["foreground"])
}

func TestImportVSCodeInclude(t *testing.T) {
	data := []byte(`{
		"name": "Child Theme",
		"type": "light",
		"include": "./base-theme.json",
		"colors": {},
		"tokenColors": []
	}`)

	theme, err := importer.ImportVSCode(data)
	require.NoError(t, err)
	assert.Equal(t, "./base-theme.json", theme.Colors["_include"])
}

func TestImportVSCodeInvalidJSON(t *testing.T) {
	_, err := importer.ImportVSCode([]byte("not json"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid VS Code theme JSON")
}

func TestImportVSCodeLightTheme(t *testing.T) {
	data := []byte(`{
		"name": "Solarized Light",
		"type": "vs",
		"colors": {},
		"tokenColors": []
	}`)

	theme, err := importer.ImportVSCode(data)
	require.NoError(t, err)
	assert.Equal(t, "light", theme.Type)
}

func TestImportVSCodeHighContrastTheme(t *testing.T) {
	data := []byte(`{
		"name": "HC Black",
		"type": "hc-black",
		"colors": {},
		"tokenColors": []
	}`)

	theme, err := importer.ImportVSCode(data)
	require.NoError(t, err)
	assert.Equal(t, "high-contrast", theme.Type)
}

func TestImportVSCodeEmptyName(t *testing.T) {
	data := []byte(`{
		"type": "dark",
		"colors": {"editor.background": "#000"},
		"tokenColors": []
	}`)

	theme, err := importer.ImportVSCode(data)
	require.NoError(t, err)
	assert.Equal(t, "Imported VS Code Theme", theme.Name)
	assert.Equal(t, "imported-vs-code-theme", theme.ID)
}

func TestImportAutoDetectsVSCode(t *testing.T) {
	theme, err := importer.Import(sampleVSCodeTheme)
	require.NoError(t, err)
	assert.Equal(t, "vscode", theme.Source)
	assert.Equal(t, "One Dark Pro", theme.Name)
}

func TestVSCodeScopeAsCommaSeparated(t *testing.T) {
	data := []byte(`{
		"name": "Scope Test",
		"type": "dark",
		"colors": {},
		"tokenColors": [
			{
				"name": "Multi",
				"scope": "keyword.control, storage.type",
				"settings": {"foreground": "#FF0000"}
			}
		]
	}`)

	theme, err := importer.ImportVSCode(data)
	require.NoError(t, err)
	require.Len(t, theme.TokenColors, 1)
	assert.Equal(t, []string{"keyword.control", "storage.type"}, theme.TokenColors[0].Scope)
}
