import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import {
  fetchThemes,
  fetchActiveTheme,
  setActiveTheme,
  fetchTheme,
} from '../api/themes';
import { useThemeStore } from '../stores/theme-store';
import type { ThemeDef } from '../types/theme';

const THEME_KEYS = {
  all: ['themes'] as const,
  list: () => [...THEME_KEYS.all, 'list'] as const,
  active: () => [...THEME_KEYS.all, 'active'] as const,
  detail: (id: string) => [...THEME_KEYS.all, 'detail', id] as const,
};

/** Fetch all available themes and sync to store. */
export function useThemes() {
  const setThemes = useThemeStore((s) => s.setThemes);

  return useQuery({
    queryKey: THEME_KEYS.list(),
    queryFn: async () => {
      const themes = await fetchThemes();
      setThemes(themes);
      return themes;
    },
    staleTime: 60_000,
  });
}

/** Fetch the active theme and sync to store. */
export function useActiveTheme() {
  const store = useThemeStore();

  return useQuery({
    queryKey: THEME_KEYS.active(),
    queryFn: async () => {
      const theme = await fetchActiveTheme();
      store.setActiveTheme(theme);
      return theme;
    },
    staleTime: 30_000,
  });
}

/** Mutation to set the active theme. Optimistically updates the store. */
export function useSetActiveTheme() {
  const queryClient = useQueryClient();
  const store = useThemeStore();

  return useMutation({
    mutationFn: (id: string) => setActiveTheme(id),

    onMutate: async (id) => {
      await queryClient.cancelQueries({ queryKey: THEME_KEYS.active() });

      const previousTheme = queryClient.getQueryData<ThemeDef>(THEME_KEYS.active());
      const allThemes = store.themes;
      const optimistic = allThemes.find((t) => t.id === id);

      if (optimistic) {
        store.setActiveTheme(optimistic);
        queryClient.setQueryData(THEME_KEYS.active(), optimistic);
      }

      return { previousTheme };
    },

    onError: (_err, _id, context) => {
      if (context?.previousTheme) {
        store.setActiveTheme(context.previousTheme);
        queryClient.setQueryData(THEME_KEYS.active(), context.previousTheme);
      }
    },

    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: THEME_KEYS.active() });
    },
  });
}

/** Fetch a single theme by ID. */
export function useTheme(id: string) {
  return useQuery({
    queryKey: THEME_KEYS.detail(id),
    queryFn: () => fetchTheme(id),
    enabled: !!id,
  });
}
