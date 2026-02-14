package types

// ThemeDef represents a complete theme definition.
type ThemeDef struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Author      string            `json:"author"`
	Type        string            `json:"type"`
	Source      string            `json:"source,omitempty"`
	Colors      map[string]string `json:"colors"`
	TokenColors []TokenColor      `json:"token_colors,omitempty"`
}

// TokenColor defines syntax highlighting colors.
type TokenColor struct {
	Name     string            `json:"name"`
	Scope    []string          `json:"scope"`
	Settings map[string]string `json:"settings"`
}

// ThemeChangeEvent is emitted when the active theme changes.
type ThemeChangeEvent struct {
	OldThemeID string `json:"old_theme_id"`
	NewThemeID string `json:"new_theme_id"`
}
