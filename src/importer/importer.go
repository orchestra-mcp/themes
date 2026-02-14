package importer

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/orchestra-mcp/themes/src/types"
)

// Format constants for theme source detection.
const (
	FormatVSCodeJSON    = "vscode-json"
	FormatTmTheme       = "tmtheme"
	FormatOrchestraJSON = "orchestra-json"
)

// DetectFormat inspects raw bytes and returns the detected theme format.
// Returns one of: "vscode-json", "tmtheme", or "orchestra-json".
func DetectFormat(data []byte) string {
	trimmed := bytes.TrimSpace(data)

	// XML plist files start with <?xml or <plist or <!DOCTYPE plist
	if bytes.HasPrefix(trimmed, []byte("<?xml")) ||
		bytes.HasPrefix(trimmed, []byte("<plist")) ||
		bytes.HasPrefix(trimmed, []byte("<!DOCTYPE plist")) {
		return FormatTmTheme
	}

	// Try to parse as JSON and inspect fields.
	if len(trimmed) > 0 && trimmed[0] == '{' {
		var probe map[string]json.RawMessage
		if err := json.Unmarshal(trimmed, &probe); err == nil {
			// Orchestra themes always have "id"; check this first since
			// both formats share "colors".
			_, hasID := probe["id"]
			if hasID {
				return FormatOrchestraJSON
			}
			// VS Code themes are identified by "tokenColors" (camelCase).
			// "colors" alone is ambiguous — Orchestra themes can also have it.
			_, hasTokenColors := probe["tokenColors"]
			if hasTokenColors {
				return FormatVSCodeJSON
			}
		}
	}

	return FormatOrchestraJSON
}

// Import auto-detects the theme format and parses accordingly.
func Import(data []byte) (*types.ThemeDef, error) {
	format := DetectFormat(data)

	switch format {
	case FormatVSCodeJSON:
		return ImportVSCode(data)
	case FormatTmTheme:
		return ImportTmTheme(data)
	case FormatOrchestraJSON:
		return importOrchestra(data)
	default:
		return nil, fmt.Errorf("unsupported theme format: %s", format)
	}
}

// importOrchestra parses an Orchestra-native JSON theme.
func importOrchestra(data []byte) (*types.ThemeDef, error) {
	var theme types.ThemeDef
	if err := json.Unmarshal(data, &theme); err != nil {
		return nil, fmt.Errorf("invalid Orchestra theme JSON: %w", err)
	}
	if theme.ID == "" {
		return nil, fmt.Errorf("theme ID is required")
	}
	if theme.Source == "" {
		theme.Source = "orchestra"
	}
	return &theme, nil
}
