import { apiClient } from './client';
import type { ThemeDef, ThemeListResponse, SetActiveThemeResponse } from '../types/theme';

const BASE = '/themes';

/** Fetch all available themes. */
export async function fetchThemes(): Promise<ThemeDef[]> {
  const { data } = await apiClient.get<ThemeListResponse>(BASE);
  return data.themes;
}

/** Fetch the currently active theme. */
export async function fetchActiveTheme(): Promise<ThemeDef> {
  const { data } = await apiClient.get<ThemeDef>(`${BASE}/active`);
  return data;
}

/** Set the active theme by ID. */
export async function setActiveTheme(id: string): Promise<SetActiveThemeResponse> {
  const { data } = await apiClient.put<SetActiveThemeResponse>(`${BASE}/active`, { id });
  return data;
}

/** Fetch a single theme by ID. */
export async function fetchTheme(id: string): Promise<ThemeDef> {
  const { data } = await apiClient.get<ThemeDef>(`${BASE}/${id}`);
  return data;
}
