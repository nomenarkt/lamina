import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';

const API_BASE = '/api/v1/admin/roles';

export function useAssignedRoles(userId?: number, orgUnitId?: number) {
  return useQuery({
    queryKey: ['roles', userId, orgUnitId],
    queryFn: async () => {
      if (!userId || !orgUnitId) return [];
      const res = await fetch(`${API_BASE}?user_id=${userId}&org_unit_id=${orgUnitId}`);
      if (!res.ok) throw new Error('Failed to fetch roles');
      return res.json();
    },
    enabled: !!userId && !!orgUnitId,
  });
}

export function useAssignRole() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (data: { user_id: number; function: string; org_unit_id: number }) => {
      const res = await fetch(API_BASE, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data),
      });
      if (!res.ok) throw new Error('Failed to assign role');
    },
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['roles'] }),
  });
}

export function useRemoveRole() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (data: { user_id: number; function: string; org_unit_id: number }) => {
      const res = await fetch(API_BASE, {
        method: 'DELETE',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data),
      });
      if (!res.ok) throw new Error('Failed to remove role');
    },
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['roles'] }),
  });
}
