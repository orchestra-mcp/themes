package providers

import (
	"fmt"

	"github.com/orchestra-mcp/framework/app/plugins"
)

// McpTools returns MCP tool definitions contributed by the Themes plugin.
func (p *ThemesPlugin) McpTools() []plugins.McpToolDefinition {
	return []plugins.McpToolDefinition{
		{
			Name:        "list_themes",
			Description: "List all available themes",
			InputSchema: map[string]any{},
			Handler:     p.toolListThemes,
		},
		{
			Name:        "get_active_theme",
			Description: "Get the currently active theme",
			InputSchema: map[string]any{},
			Handler:     p.toolGetActiveTheme,
		},
		{
			Name:        "set_active_theme",
			Description: "Switch the active theme by ID",
			InputSchema: map[string]any{
				"id": map[string]any{
					"type":        "string",
					"description": "Theme ID to activate",
				},
			},
			Handler: p.toolSetActiveTheme,
		},
	}
}

func (p *ThemesPlugin) toolListThemes(_ map[string]any) (any, error) {
	themes := p.svc.GetAvailableThemes()
	return map[string]any{"themes": themes}, nil
}

func (p *ThemesPlugin) toolGetActiveTheme(_ map[string]any) (any, error) {
	theme := p.svc.GetActiveTheme()
	if theme == nil {
		return nil, fmt.Errorf("no active theme")
	}
	return theme, nil
}

func (p *ThemesPlugin) toolSetActiveTheme(input map[string]any) (any, error) {
	id, _ := input["id"].(string)
	if id == "" {
		return nil, fmt.Errorf("theme id is required")
	}
	if err := p.svc.SetActiveTheme(id); err != nil {
		return nil, err
	}
	return map[string]any{"active_theme": id}, nil
}
