package providers

import (
	"github.com/orchestra-mcp/framework/app/plugins"
	"github.com/orchestra-mcp/themes/config"
	"github.com/orchestra-mcp/themes/src/service"
)

// ThemesPlugin implements the Orchestra plugin interface for themes.
type ThemesPlugin struct {
	active bool
	ctx    *plugins.PluginContext
	cfg    *config.ThemesConfig
	svc    *service.ThemesService
}

// NewThemesPlugin creates a new Themes plugin instance.
func NewThemesPlugin() *ThemesPlugin { return &ThemesPlugin{} }

func (p *ThemesPlugin) ID() string             { return "orchestra/themes" }
func (p *ThemesPlugin) Name() string           { return "Themes" }
func (p *ThemesPlugin) Version() string        { return "0.1.0" }
func (p *ThemesPlugin) Dependencies() []string { return []string{"orchestra/settings"} }
func (p *ThemesPlugin) IsActive() bool         { return p.active }
func (p *ThemesPlugin) FeatureFlag() string    { return "themes" }
func (p *ThemesPlugin) ConfigKey() string      { return "themes" }

func (p *ThemesPlugin) DefaultConfig() map[string]any {
	return map[string]any{"default_theme": "orchestra-dark"}
}

// Activate initializes the themes service with built-in themes.
func (p *ThemesPlugin) Activate(ctx *plugins.PluginContext) error {
	p.ctx = ctx
	p.cfg = config.DefaultConfig()

	if dt := ctx.GetConfigString("default_theme"); dt != "" {
		p.cfg.DefaultTheme = dt
	}

	p.svc = service.New(ctx.StoragePath, p.cfg.DefaultTheme, ctx.Logger)
	p.active = true
	ctx.Logger.Info().Str("plugin", p.ID()).Msg("themes plugin activated")
	return nil
}

// Deactivate shuts down the themes plugin.
func (p *ThemesPlugin) Deactivate() error {
	p.active = false
	return nil
}

// Service returns the underlying ThemesService.
func (p *ThemesPlugin) Service() *service.ThemesService {
	return p.svc
}

// Contributes returns theme contributions for the ContributesRegistry.
func (p *ThemesPlugin) Contributes() *plugins.Contributions {
	return &plugins.Contributions{
		Themes: []plugins.ThemeContribution{
			{
				ID:          "orchestra-light",
				Label:       "Orchestra Light",
				UITheme:     "light",
				Description: "Default light theme for Orchestra",
			},
			{
				ID:          "orchestra-dark",
				Label:       "Orchestra Dark",
				UITheme:     "dark",
				Description: "Default dark theme for Orchestra",
			},
		},
	}
}
