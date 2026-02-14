# Orchestra Themes Plugin

Color theme management for Orchestra. Ships with 2 built-in themes (light/dark), supports VS Code JSON and `.tmTheme` XML import, persists user preference.

## Features

- **Built-in themes** — Orchestra Light and Orchestra Dark with 24 color tokens each
- **VS Code import** — import VS Code JSON themes and `.tmTheme` XML (TextMate)
- **Theme switching** — change active theme with listener notifications
- **Export/import** — serialize themes to JSON for sharing
- **Preference persistence** — saves active theme to `theme-preference.json`

## Configuration

| Field | Default | Description |
|-------|---------|-------------|
| `DefaultTheme` | `orchestra-dark` | Theme ID activated on first launch |

## MCP Tools

| Tool | Description |
|------|-------------|
| `list_themes` | All available themes |
| `get_active_theme` | Currently active theme |
| `set_active_theme` | Switch theme by ID |

## REST API

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/themes/` | List all themes |
| `GET` | `/themes/active` | Get active theme |
| `PUT` | `/themes/active` | Set active theme |
| `GET` | `/themes/:id` | Get specific theme |
| `POST` | `/themes/import` | Import custom theme (JSON) |
| `POST` | `/themes/import/vscode` | Import VS Code/tmTheme format |
| `GET` | `/themes/:id/export` | Export theme as JSON |

## Package Structure

```
plugins/themes/
├── config/themes.go           # ThemesConfig
├── providers/
│   ├── plugin.go              # ThemesPlugin (activate, services, tools)
│   ├── routes.go              # REST endpoints + import handlers
│   └── tools.go               # 3 MCP tool definitions
├── src/
│   ├── builtin/themes.go      # Light and Dark theme definitions
│   ├── importer/
│   │   ├── importer.go        # Format detection + unified Import()
│   │   ├── vscode.go          # VS Code JSON import + color mapping
│   │   ├── tmtheme.go         # tmTheme XML import + color extraction
│   │   └── plist.go           # Plist XML decoder (types + parser)
│   ├── service/service.go     # ThemesService (register, activate, export)
│   └── types/types.go         # ThemeDef, TokenColor, ThemeChangeEvent
├── tests/
│   ├── service_test.go        # Theme registration, activation, persistence
│   ├── importer_test.go       # Format detection + VS Code import tests
│   └── tmtheme_test.go        # tmTheme import + unified import + slugify
└── go.mod
```
