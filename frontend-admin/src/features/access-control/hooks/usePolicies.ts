import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';

const API_BASE = '/api/v1/admin/policies';

export function usePolicies() {
  return useQuery({
    queryKey: ['policies'],
    queryFn: async () => {
      const res = await fetch(API_BASE);
      if (!res.ok) throw new Error('Failed to load policies');
      return res.json(); // Expected shape: [["planner", "orgunit:47", "/api/flights", "write"], ...]
    },
  });
}

export function useAddPolicy() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (data: {
      role: string;
      org_unit_id: number;
      object: string;
      action: string;
    }) => {
      const res = await fetch(API_BASE, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data),
      });
      if (!res.ok) throw new Error('Failed to add policy');
    },
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['policies'] }),
  });
}

export function useDeletePolicy() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (data: {
      role: string;
      org_unit_id: number;
      object: string;
      action: string;
    }) => {
      const res = await fetch(API_BASE, {
        method: 'DELETE',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data),
      });
      if (!res.ok) throw new Error('Failed to delete policy');
    },
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['policies'] }),
  });
}
