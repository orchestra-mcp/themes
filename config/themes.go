package config

// ThemesConfig holds configuration for the Themes plugin.
type ThemesConfig struct {
	DefaultTheme string `json:"default_theme"`
}

// DefaultConfig returns the default themes configuration.
func DefaultConfig() *ThemesConfig {
	return &ThemesConfig{DefaultTheme: "orchestra-dark"}
}
