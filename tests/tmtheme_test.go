package tests

import (
	"testing"

	"github.com/orchestra-mcp/themes/src/importer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- TextMate .tmTheme Import ---

var sampleTmTheme = []byte(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>name</key>
	<string>Monokai</string>
	<key>author</key>
	<string>Wimer Hazenberg</string>
	<key>settings</key>
	<array>
		<dict>
			<key>settings</key>
			<dict>
				<key>background</key>
				<string>#272822</string>
				<key>foreground</key>
				<string>#F8F8F2</string>
				<key>caret</key>
				<string>#F8F8F0</string>
				<key>selection</key>
				<string>#49483E</string>
				<key>lineHighlight</key>
				<string>#3E3D32</string>
			</dict>
		</dict>
		<dict>
			<key>name</key>
			<string>Comment</string>
			<key>scope</key>
			<string>comment</string>
			<key>settings</key>
			<dict>
				<key>foreground</key>
				<string>#75715E</string>
			</dict>
		</dict>
		<dict>
			<key>name</key>
			<string>String</string>
			<key>scope</key>
			<string>string.quoted, string.interpolated</string>
			<key>settings</key>
			<dict>
				<key>foreground</key>
				<string>#E6DB74</string>
			</dict>
		</dict>
		<dict>
			<key>name</key>
			<string>Keyword</string>
			<key>scope</key>
			<string>keyword</string>
			<key>settings</key>
			<dict>
				<key>foreground</key>
				<string>#F92672</string>
				<key>fontStyle</key>
				<string>bold</string>
			</dict>
		</dict>
	</array>
</dict>
</plist>`)

func TestImportTmThemeBasicFields(t *testing.T) {
	theme, err := importer.ImportTmTheme(sampleTmTheme)
	require.NoError(t, err)

	assert.Equal(t, "monokai", theme.ID)
	assert.Equal(t, "Monokai", theme.Name)
	assert.Equal(t, "Wimer Hazenberg", theme.Author)
	assert.Equal(t, "tmtheme", theme.Source)
	assert.Equal(t, "dark", theme.Type)
}

func TestImportTmThemeGlobalColors(t *testing.T) {
	theme, err := importer.ImportTmTheme(sampleTmTheme)
	require.NoError(t, err)

	assert.Equal(t, "#272822", theme.Colors["bg-primary"])
	assert.Equal(t, "#F8F8F2", theme.Colors["text-primary"])
	assert.Equal(t, "#F8F8F0", theme.Colors["caret"])
	assert.Equal(t, "#49483E", theme.Colors["bg-selection"])
	assert.Equal(t, "#3E3D32", theme.Colors["bg-line-highlight"])
	assert.Equal(t, "#272822", theme.Colors["raw.background"])
	assert.Equal(t, "#F8F8F2", theme.Colors["raw.foreground"])
}

func TestImportTmThemeTokenColors(t *testing.T) {
	theme, err := importer.ImportTmTheme(sampleTmTheme)
	require.NoError(t, err)

	require.Len(t, theme.TokenColors, 3)

	assert.Equal(t, "Comment", theme.TokenColors[0].Name)
	assert.Equal(t, []string{"comment"}, theme.TokenColors[0].Scope)
	assert.Equal(t, "#75715E", theme.TokenColors[0].Settings["foreground"])

	assert.Equal(t, "String", theme.TokenColors[1].Name)
	assert.Contains(t, theme.TokenColors[1].Scope, "string.quoted")
	assert.Contains(t, theme.TokenColors[1].Scope, "string.interpolated")
	assert.Equal(t, "#E6DB74", theme.TokenColors[1].Settings["foreground"])

	assert.Equal(t, "Keyword", theme.TokenColors[2].Name)
	assert.Equal(t, "#F92672", theme.TokenColors[2].Settings["foreground"])
	assert.Equal(t, "bold", theme.TokenColors[2].Settings["fontStyle"])
}

func TestImportTmThemeInvalidXML(t *testing.T) {
	_, err := importer.ImportTmTheme([]byte("not xml"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid tmTheme XML")
}

func TestImportTmThemeNoSettings(t *testing.T) {
	data := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<plist version="1.0">
<dict>
	<key>name</key>
	<string>Empty</string>
</dict>
</plist>`)

	_, err := importer.ImportTmTheme(data)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no settings array")
}

func TestImportTmThemeLightDetection(t *testing.T) {
	data := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<plist version="1.0">
<dict>
	<key>name</key>
	<string>Light Theme</string>
	<key>settings</key>
	<array>
		<dict>
			<key>settings</key>
			<dict>
				<key>background</key>
				<string>#FAFAFA</string>
				<key>foreground</key>
				<string>#333333</string>
			</dict>
		</dict>
	</array>
</dict>
</plist>`)

	theme, err := importer.ImportTmTheme(data)
	require.NoError(t, err)
	assert.Equal(t, "light", theme.Type)
}

// --- Unified Import ---

func TestImportAutoDetectsTmTheme(t *testing.T) {
	theme, err := importer.Import(sampleTmTheme)
	require.NoError(t, err)
	assert.Equal(t, "tmtheme", theme.Source)
	assert.Equal(t, "Monokai", theme.Name)
}

func TestImportAutoDetectsOrchestra(t *testing.T) {
	data := []byte(`{
		"id": "custom",
		"name": "Custom",
		"type": "dark",
		"colors": {"background": "#000"}
	}`)

	theme, err := importer.Import(data)
	require.NoError(t, err)
	assert.Equal(t, "orchestra", theme.Source)
	assert.Equal(t, "custom", theme.ID)
}

func TestImportOrchestraMissingID(t *testing.T) {
	data := []byte(`{"name": "No ID", "colors": {}}`)
	_, err := importer.Import(data)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "theme ID is required")
}

func TestImportInvalidData(t *testing.T) {
	_, err := importer.Import([]byte("totally broken"))
	assert.Error(t, err)
}

// --- Slugify ---

func TestSlugifyVariousNames(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"One Dark Pro", "one-dark-pro"},
		{"Monokai++", "monokai"},
		{"Solarized (Light)", "solarized-light"},
		{"  Spacey Theme  ", "spacey-theme"},
		{"UPPER_case_Mix", "upper-case-mix"},
	}

	for _, tc := range tests {
		data := []byte(`{"name":"` + tc.input + `","type":"dark","colors":{},"tokenColors":[]}`)
		theme, err := importer.ImportVSCode(data)
		require.NoError(t, err, "input: %s", tc.input)
		assert.Equal(t, tc.expected, theme.ID, "input: %s", tc.input)
	}
}
