package service

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/orchestra-mcp/themes/src/builtin"
	"github.com/orchestra-mcp/themes/src/types"
	"github.com/rs/zerolog"
)

// ThemesService manages theme registration, activation, and persistence.
type ThemesService struct {
	mu          sync.RWMutex
	themes      map[string]*types.ThemeDef
	activeID    string
	storagePath string
	listeners   []func(old, new *types.ThemeDef)
	logger      zerolog.Logger
}

// New creates a ThemesService with built-in themes loaded.
func New(storagePath, defaultTheme string, logger zerolog.Logger) *ThemesService {
	svc := &ThemesService{
		themes:      make(map[string]*types.ThemeDef),
		activeID:    defaultTheme,
		storagePath: storagePath,
		logger:      logger,
	}

	for _, t := range builtin.BuiltinThemes() {
		svc.themes[t.ID] = t
	}

	svc.loadPreference()
	return svc
}

// RegisterTheme registers a theme definition.
func (s *ThemesService) RegisterTheme(theme *types.ThemeDef) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.themes[theme.ID] = theme
	s.logger.Info().Str("theme", theme.ID).Msg("theme registered")
}

// SetActiveTheme switches to a theme by ID.
func (s *ThemesService) SetActiveTheme(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	newTheme, ok := s.themes[id]
	if !ok {
		return fmt.Errorf("theme not found: %s", id)
	}

	oldID := s.activeID
	oldTheme := s.themes[oldID]
	s.activeID = id
	s.savePreference()

	s.mu.Unlock()
	s.fireListeners(oldTheme, newTheme)
	s.mu.Lock()

	return nil
}

// GetActiveTheme returns the currently active theme.
func (s *ThemesService) GetActiveTheme() *types.ThemeDef {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.themes[s.activeID]
}

// GetAvailableThemes returns all registered themes.
func (s *ThemesService) GetAvailableThemes() []types.ThemeDef {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]types.ThemeDef, 0, len(s.themes))
	for _, t := range s.themes {
		result = append(result, *t)
	}
	return result
}

// GetTheme returns a specific theme by ID.
func (s *ThemesService) GetTheme(id string) (*types.ThemeDef, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	t, ok := s.themes[id]
	if !ok {
		return nil, fmt.Errorf("theme not found: %s", id)
	}
	return t, nil
}

// OnDidChangeTheme registers a callback for theme changes.
func (s *ThemesService) OnDidChangeTheme(cb func(old, new *types.ThemeDef)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.listeners = append(s.listeners, cb)
}

// ExportTheme serializes a theme to JSON bytes.
func (s *ThemesService) ExportTheme(id string) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	t, ok := s.themes[id]
	if !ok {
		return nil, fmt.Errorf("theme not found: %s", id)
	}
	return json.MarshalIndent(t, "", "  ")
}

// ImportTheme deserializes a theme from JSON and registers it.
func (s *ThemesService) ImportTheme(data []byte) (*types.ThemeDef, error) {
	var theme types.ThemeDef
	if err := json.Unmarshal(data, &theme); err != nil {
		return nil, fmt.Errorf("invalid theme JSON: %w", err)
	}
	if theme.ID == "" {
		return nil, fmt.Errorf("theme ID is required")
	}
	s.RegisterTheme(&theme)
	return &theme, nil
}

func (s *ThemesService) fireListeners(old, new *types.ThemeDef) {
	for _, cb := range s.listeners {
		cb(old, new)
	}
}

type preference struct {
	ActiveTheme string `json:"active_theme"`
}

func (s *ThemesService) prefPath() string {
	return filepath.Join(s.storagePath, "theme-preference.json")
}

func (s *ThemesService) loadPreference() {
	data, err := os.ReadFile(s.prefPath())
	if err != nil {
		return
	}
	var pref preference
	if err := json.Unmarshal(data, &pref); err != nil {
		return
	}
	if _, ok := s.themes[pref.ActiveTheme]; ok {
		s.activeID = pref.ActiveTheme
	}
}

func (s *ThemesService) savePreference() {
	data, err := json.Marshal(preference{ActiveTheme: s.activeID})
	if err != nil {
		s.logger.Warn().Err(err).Msg("failed to marshal theme preference")
		return
	}
	if err := os.MkdirAll(s.storagePath, 0o755); err != nil {
		s.logger.Warn().Err(err).Msg("failed to create storage dir")
		return
	}
	if err := os.WriteFile(s.prefPath(), data, 0o644); err != nil {
		s.logger.Warn().Err(err).Msg("failed to save theme preference")
	}
}
