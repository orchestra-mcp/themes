package providers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/orchestra-mcp/themes/src/importer"
	"github.com/orchestra-mcp/themes/src/types"
)

// RegisterRoutes registers all REST API routes for the themes plugin.
func (p *ThemesPlugin) RegisterRoutes(group fiber.Router) {
	themes := group.Group("/themes")

	themes.Get("/", p.handleListThemes)
	themes.Get("/active", p.handleGetActive)
	themes.Put("/active", p.handleSetActive)
	themes.Get("/:id", p.handleGetTheme)
	themes.Post("/import", p.handleImport)
	themes.Post("/import/vscode", p.handleImportVSCode)
	themes.Get("/:id/export", p.handleExport)
}

func (p *ThemesPlugin) handleListThemes(c fiber.Ctx) error {
	themes := p.svc.GetAvailableThemes()
	return c.JSON(fiber.Map{"themes": themes})
}

func (p *ThemesPlugin) handleGetActive(c fiber.Ctx) error {
	theme := p.svc.GetActiveTheme()
	if theme == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "not_found",
			"message": "No active theme set",
		})
	}
	return c.JSON(theme)
}

type setActiveRequest struct {
	ID string `json:"id"`
}

func (p *ThemesPlugin) handleSetActive(c fiber.Ctx) error {
	var req setActiveRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "Invalid JSON body",
		})
	}
	if req.ID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "validation_error",
			"message": "Theme ID is required",
		})
	}
	if err := p.svc.SetActiveTheme(req.ID); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "not_found",
			"message": err.Error(),
		})
	}
	return c.JSON(fiber.Map{"active_theme": req.ID})
}

func (p *ThemesPlugin) handleGetTheme(c fiber.Ctx) error {
	id := c.Params("id")
	theme, err := p.svc.GetTheme(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "not_found",
			"message": err.Error(),
		})
	}
	return c.JSON(theme)
}

func (p *ThemesPlugin) handleImport(c fiber.Ctx) error {
	body := c.Body()
	if len(body) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "Request body is empty",
		})
	}
	theme, err := p.svc.ImportTheme(body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "import_error",
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(theme)
}

func (p *ThemesPlugin) handleImportVSCode(c fiber.Ctx) error {
	body := c.Body()
	if len(body) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "Request body is empty",
		})
	}

	format := importer.DetectFormat(body)
	var theme *types.ThemeDef
	var err error

	switch format {
	case importer.FormatVSCodeJSON:
		theme, err = importer.ImportVSCode(body)
	case importer.FormatTmTheme:
		theme, err = importer.ImportTmTheme(body)
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_format",
			"message": "Expected VS Code JSON theme or .tmTheme XML",
		})
	}

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "import_error",
			"message": err.Error(),
		})
	}

	p.svc.RegisterTheme(theme)
	return c.Status(fiber.StatusCreated).JSON(theme)
}

func (p *ThemesPlugin) handleExport(c fiber.Ctx) error {
	id := c.Params("id")
	data, err := p.svc.ExportTheme(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "not_found",
			"message": err.Error(),
		})
	}
	c.Set("Content-Type", "application/json")
	return c.Send(data)
}
