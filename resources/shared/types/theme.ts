/**
 * Theme type definitions matching the Go ThemeDef struct
 * from plugins/themes/src/types/types.go.
 */

/** Syntax highlighting token color definition. */
export interface TokenColor {
  name: string;
  scope: string[];
  settings: Record<string, string>;
}

/** Complete theme definition returned by the themes API. */
export interface ThemeDef {
  id: string;
  name: string;
  description: string;
  author: string;
  type: 'light' | 'dark';
  source?: string;
  colors: Record<string, string>;
  token_colors?: TokenColor[];
}

/** Event emitted when the active theme changes. */
export interface ThemeChangeEvent {
  old_theme_id: string;
  new_theme_id: string;
}

/** API response shape for listing themes. */
export interface ThemeListResponse {
  themes: ThemeDef[];
}

/** API response shape for setting the active theme. */
export interface SetActiveThemeResponse {
  active_theme: string;
}
