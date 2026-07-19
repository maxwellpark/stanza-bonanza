import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '@/lib/api';
import type { Poem, CreatePoemInput, SubmitStanzaInput, Stanza } from '@/types/poem';
import type { PaginatedResponse, PaginationParams } from '@/types/api';

export function usePoems(params?: PaginationParams & { format?: string; sort?: string; q?: string; tag?: string }) {
  var query = new URLSearchParams();
  if (params?.page) { query.set('page', String(params.page)); }
  if (params?.pageSize) { query.set('pageSize', String(params.pageSize)); }
  if (params?.format) { query.set('format', params.format); }
  if (params?.sort) { query.set('sort', params.sort); }
  if (params?.q) { query.set('q', params.q); }
  if (params?.tag) { query.set('tag', params.tag); }

  return useQuery({
    queryKey: ['poems', params],
    queryFn: () => api.get<PaginatedResponse<Poem>>(`/poems?${query}`),
  });
}

export function usePoem(id: string) {
  return useQuery({
    queryKey: ['poem', id],
    queryFn: () => api.get<Poem>(`/poems/${id}`),
    enabled: !!id,
  });
}

export function useDrafts() {
  return useQuery({
    queryKey: ['drafts'],
    queryFn: () => api.get<PaginatedResponse<Poem>>('/poems/drafts'),
  });
}

export function usePublishPoem(poemId: string) {
  var queryClient = useQueryClient();
  return useMutation({
    mutationFn: () => api.post(`/poems/${poemId}/publish`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['poem', poemId] });
      queryClient.invalidateQueries({ queryKey: ['drafts'] });
      queryClient.invalidateQueries({ queryKey: ['poems'] });
    },
  });
}

export function useCreatePoem() {
  var queryClient = useQueryClient();
  return useMutation({
    mutationFn: (input: CreatePoemInput) => api.post<Poem>('/poems', input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['poems'] });
    },
  });
}

export function useSubmitStanza(poemId: string) {
  var queryClient = useQueryClient();
  return useMutation({
    mutationFn: (input: SubmitStanzaInput) =>
      api.post<Stanza>(`/poems/${poemId}/stanzas`, input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['poem', poemId] });
    },
  });
}

export function useExplore(params?: PaginationParams) {
  return useQuery({
    queryKey: ['explore', params],
    queryFn: () => api.get<PaginatedResponse<Poem>>(`/explore?page=${params?.page ?? 1}`),
  });
}

export function useFeed(params?: PaginationParams) {
  return useQuery({
    queryKey: ['feed', params],
    queryFn: () => api.get<PaginatedResponse<Poem>>(`/feed?page=${params?.page ?? 1}`),
  });
}

export function useHallOfFame(params?: PaginationParams) {
  return useQuery({
    queryKey: ['hallOfFame', params],
    queryFn: () => api.get<PaginatedResponse<Poem>>(`/hall-of-fame?page=${params?.page ?? 1}`),
  });
}

export function useUserPoems(userId: string, params?: PaginationParams) {
  return useQuery({
    queryKey: ['userPoems', userId, params],
    queryFn: () => api.get<PaginatedResponse<Poem>>(`/users/${userId}/poems?page=${params?.page ?? 1}`),
    enabled: !!userId,
  });
}

export function useReviewStanza(poemId: string) {
  var queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ stanzaId, approved }: { stanzaId: string; approved: boolean }) =>
      api.put(`/poems/${poemId}/stanzas/${stanzaId}`, { approved }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['poem', poemId] });
    },
  });
}
