import { useState } from 'react';
import type { FC } from 'react';
import { Check, ChevronDown, Palette } from 'lucide-react';
import { useThemeContext } from '@orchestra/shared/providers/theme-provider';
import type { ThemeDef } from '@orchestra/shared/types/theme';
import { cn } from '../../lib/utils';

// ---------------------------------------------------------------------------
// Color swatch preview
// ---------------------------------------------------------------------------

const SWATCH_KEYS = ['background', 'primary', 'accent', 'foreground'] as const;

interface SwatchProps {
  theme: ThemeDef;
}

const ThemeSwatches: FC<SwatchProps> = ({ theme }) => (
  <div className="flex items-center gap-0.5">
    {SWATCH_KEYS.map((key) => {
      const color = theme.colors[key];
      if (!color) return null;
      return (
        <span
          key={key}
          className="inline-block size-3 rounded-full border border-gray-300 dark:border-gray-600"
          style={{ backgroundColor: color }}
          aria-hidden
        />
      );
    })}
  </div>
);

// ---------------------------------------------------------------------------
// Theme row
// ---------------------------------------------------------------------------

interface ThemeRowProps {
  theme: ThemeDef;
  isActive: boolean;
  onSelect: (id: string) => void;
}

const ThemeRow: FC<ThemeRowProps> = ({ theme, isActive, onSelect }) => (
  <button
    type="button"
    className={cn(
      'flex w-full items-center gap-3 rounded-md px-3 py-2 text-sm',
      'hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors',
      isActive && 'bg-gray-100 dark:bg-gray-800',
    )}
    onClick={() => onSelect(theme.id)}
  >
    <ThemeSwatches theme={theme} />
    <span className="flex-1 text-left truncate">{theme.name}</span>
    {isActive && <Check className="size-4 text-green-600 dark:text-green-400 shrink-0" />}
  </button>
);

// ---------------------------------------------------------------------------
// Theme group
// ---------------------------------------------------------------------------

interface ThemeGroupProps {
  label: string;
  themes: ThemeDef[];
  activeId: string | null;
  onSelect: (id: string) => void;
}

const ThemeGroup: FC<ThemeGroupProps> = ({ label, themes, activeId, onSelect }) => {
  if (themes.length === 0) return null;
  return (
    <div>
      <div className="px-3 py-1.5 text-xs font-semibold uppercase tracking-wider text-gray-500">
        {label}
      </div>
      {themes.map((t) => (
        <ThemeRow
          key={t.id}
          theme={t}
          isActive={t.id === activeId}
          onSelect={onSelect}
        />
      ))}
    </div>
  );
};

// ---------------------------------------------------------------------------
// ThemeSelector
// ---------------------------------------------------------------------------

export interface ThemeSelectorProps {
  /** Additional classes on the container. */
  className?: string;
}

export const ThemeSelector: FC<ThemeSelectorProps> = ({ className }) => {
  const { theme, themes, setTheme } = useThemeContext();
  const [open, setOpen] = useState(false);

  const lightThemes = themes.filter((t) => t.type === 'light');
  const darkThemes = themes.filter((t) => t.type === 'dark');

  return (
    <div className={cn('relative inline-block', className)}>
      <button
        type="button"
        className={cn(
          'flex items-center gap-2 rounded-md border px-3 py-2 text-sm',
          'border-gray-200 dark:border-gray-700',
          'hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors',
        )}
        onClick={() => setOpen((prev) => !prev)}
        aria-haspopup="listbox"
        aria-expanded={open}
      >
        <Palette className="size-4" />
        <span className="truncate max-w-[140px]">{theme?.name ?? 'Select theme'}</span>
        <ChevronDown className={cn('size-4 transition-transform', open && 'rotate-180')} />
      </button>

      {open && (
        <>
          {/* Backdrop to close the dropdown */}
          <div
            className="fixed inset-0 z-40"
            onClick={() => setOpen(false)}
            aria-hidden
          />
          <div
            className={cn(
              'absolute right-0 z-50 mt-1 w-64 rounded-lg border p-1.5 shadow-lg',
              'border-gray-200 bg-white dark:border-gray-700 dark:bg-gray-900',
            )}
            role="listbox"
            aria-label="Select a theme"
          >
            <ThemeGroup
              label="Light"
              themes={lightThemes}
              activeId={theme?.id ?? null}
              onSelect={(id) => {
                setTheme(id);
                setOpen(false);
              }}
            />
            {lightThemes.length > 0 && darkThemes.length > 0 && (
              <div className="my-1 border-t border-gray-200 dark:border-gray-700" />
            )}
            <ThemeGroup
              label="Dark"
              themes={darkThemes}
              activeId={theme?.id ?? null}
              onSelect={(id) => {
                setTheme(id);
                setOpen(false);
              }}
            />
          </div>
        </>
      )}
    </div>
  );
};
