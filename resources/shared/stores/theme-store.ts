import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import type { ThemeDef } from '../types/theme';

/** Theme store state (data only). */
export interface ThemeStoreState {
  activeTheme: ThemeDef | null;
  activeThemeId: string | null;
  themes: ThemeDef[];
  loading: boolean;
}

/** Theme store actions (behavior only). */
export interface ThemeStoreActions {
  setActiveTheme: (theme: ThemeDef) => void;
  setActiveThemeId: (id: string) => void;
  setThemes: (themes: ThemeDef[]) => void;
  setLoading: (loading: boolean) => void;
  reset: () => void;
}

const initialState: ThemeStoreState = {
  activeTheme: null,
  activeThemeId: null,
  themes: [],
  loading: false,
};

/**
 * Zustand store for theme state.
 * Persists activeThemeId to localStorage so the last selection
 * survives page reloads even before the API responds.
 */
export const useThemeStore = create<ThemeStoreState & ThemeStoreActions>()(
  persist(
    (set) => ({
      ...initialState,

      setActiveTheme: (theme) =>
        set({ activeTheme: theme, activeThemeId: theme.id }),

      setActiveThemeId: (id) =>
        set({ activeThemeId: id }),

      setThemes: (themes) =>
        set({ themes }),

      setLoading: (loading) =>
        set({ loading }),

      reset: () =>
        set(initialState),
    }),
    {
      name: 'orchestra-theme',
      partialize: (state) => ({ activeThemeId: state.activeThemeId }),
    },
  ),
);
