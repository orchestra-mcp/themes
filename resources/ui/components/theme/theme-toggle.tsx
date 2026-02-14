import type { FC } from 'react';
import { Moon, Sun } from 'lucide-react';
import { useThemeContext } from '@orchestra/shared/providers/theme-provider';
import { cn } from '../../lib/utils';

export interface ThemeToggleProps {
  /** Additional classes on the button. */
  className?: string;
  /** Size of the icon in pixels. Defaults to 18. */
  iconSize?: number;
}

/**
 * Simple light/dark toggle button.
 * Switches between the default light and dark themes using the ThemeProvider.
 */
export const ThemeToggle: FC<ThemeToggleProps> = ({
  className,
  iconSize = 18,
}) => {
  const { theme, toggleTheme, loading } = useThemeContext();
  const isDark = theme?.type === 'dark';

  return (
    <button
      type="button"
      className={cn(
        'inline-flex items-center justify-center rounded-md p-2',
        'border border-gray-200 dark:border-gray-700',
        'hover:bg-gray-100 dark:hover:bg-gray-800',
        'transition-colors focus-visible:outline-none focus-visible:ring-2',
        'focus-visible:ring-blue-500 focus-visible:ring-offset-2',
        'disabled:pointer-events-none disabled:opacity-50',
        className,
      )}
      onClick={toggleTheme}
      disabled={loading}
      aria-label={isDark ? 'Switch to light theme' : 'Switch to dark theme'}
      title={isDark ? 'Switch to light theme' : 'Switch to dark theme'}
    >
      {isDark ? (
        <Sun size={iconSize} className="text-amber-400" />
      ) : (
        <Moon size={iconSize} className="text-slate-700" />
      )}
    </button>
  );
};
