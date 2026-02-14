import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useRef,
} from 'react';
import type { FC, ReactNode } from 'react';
import { useThemeStore } from '../stores/theme-store';
import { useActiveTheme, useThemes, useSetActiveTheme } from '../hooks/use-theme';
import type { ThemeDef } from '../types/theme';

// ---------------------------------------------------------------------------
// CSS variable injection
// ---------------------------------------------------------------------------

/**
 * Convert a theme color key to a CSS custom property name.
 * Dot-notation keys like "sidebar.background" become "--sidebar-background".
 * Simple keys like "primary" become "--primary".
 */
function colorKeyToCssVar(key: string): string {
  return `--${key.replace(/\./g, '-')}`;
}

/** Inject all theme colors as CSS custom properties on an element. */
function injectCssVariables(colors: Record<string, string>, el: HTMLElement): void {
  for (const [key, value] of Object.entries(colors)) {
    el.style.setProperty(colorKeyToCssVar(key), value);
  }
}

/** Remove all previously injected theme CSS properties from an element. */
function removeCssVariables(colors: Record<string, string>, el: HTMLElement): void {
  for (const key of Object.keys(colors)) {
    el.style.removeProperty(colorKeyToCssVar(key));
  }
}

// ---------------------------------------------------------------------------
// Context
// ---------------------------------------------------------------------------

export interface ThemeContextValue {
  /** The currently active theme definition, or null while loading. */
  theme: ThemeDef | null;
  /** All available themes. */
  themes: ThemeDef[];
  /** Whether the initial load is in progress. */
  loading: boolean;
  /** Switch to a specific theme by ID. */
  setTheme: (id: string) => void;
  /** Toggle between the default light and dark theme. */
  toggleTheme: () => void;
}

const ThemeContext = createContext<ThemeContextValue | null>(null);

// ---------------------------------------------------------------------------
// Provider
// ---------------------------------------------------------------------------

export interface ThemeProviderProps {
  children: ReactNode;
  /** ID of the default light theme. Defaults to "orchestra-light". */
  defaultLightId?: string;
  /** ID of the default dark theme. Defaults to "orchestra-dark". */
  defaultDarkId?: string;
}

export const ThemeProvider: FC<ThemeProviderProps> = ({
  children,
  defaultLightId = 'orchestra-light',
  defaultDarkId = 'orchestra-dark',
}) => {
  const activeTheme = useThemeStore((s) => s.activeTheme);
  const themes = useThemeStore((s) => s.themes);
  const prevColorsRef = useRef<Record<string, string> | null>(null);

  // Fetch data from API on mount.
  const { isLoading: themesLoading } = useThemes();
  const { isLoading: activeLoading } = useActiveTheme();
  const { mutate: setActiveMutation } = useSetActiveTheme();

  const loading = themesLoading || activeLoading;

  // Detect system color scheme preference as fallback.
  useEffect(() => {
    if (activeTheme || loading) return;

    const mq = window.matchMedia('(prefers-color-scheme: dark)');
    const fallbackId = mq.matches ? defaultDarkId : defaultLightId;
    const fallback = themes.find((t) => t.id === fallbackId);
    if (fallback) {
      useThemeStore.getState().setActiveTheme(fallback);
    }
  }, [activeTheme, themes, loading, defaultDarkId, defaultLightId]);

  // Inject / update CSS variables whenever the active theme changes.
  useEffect(() => {
    if (!activeTheme) return;

    const el = document.documentElement;

    // Remove old variables first to avoid stale properties.
    if (prevColorsRef.current) {
      removeCssVariables(prevColorsRef.current, el);
    }

    injectCssVariables(activeTheme.colors, el);
    prevColorsRef.current = activeTheme.colors;

    // Set data-theme attribute for CSS selectors.
    el.setAttribute('data-theme', activeTheme.type);

    return () => {
      if (prevColorsRef.current) {
        removeCssVariables(prevColorsRef.current, el);
        prevColorsRef.current = null;
      }
    };
  }, [activeTheme]);

  const setTheme = useCallback(
    (id: string) => setActiveMutation(id),
    [setActiveMutation],
  );

  const toggleTheme = useCallback(() => {
    const currentType = activeTheme?.type ?? 'light';
    const targetId = currentType === 'light' ? defaultDarkId : defaultLightId;
    setActiveMutation(targetId);
  }, [activeTheme, defaultDarkId, defaultLightId, setActiveMutation]);

  const value = useMemo<ThemeContextValue>(
    () => ({
      theme: activeTheme,
      themes,
      loading,
      setTheme,
      toggleTheme,
    }),
    [activeTheme, themes, loading, setTheme, toggleTheme],
  );

  return (
    <ThemeContext.Provider value={value}>
      {children}
    </ThemeContext.Provider>
  );
};

// ---------------------------------------------------------------------------
// Hook
// ---------------------------------------------------------------------------

/** Access the current theme context. Must be used inside ThemeProvider. */
export function useThemeContext(): ThemeContextValue {
  const ctx = useContext(ThemeContext);
  if (!ctx) {
    throw new Error('useThemeContext must be used within a <ThemeProvider>');
  }
  return ctx;
}
