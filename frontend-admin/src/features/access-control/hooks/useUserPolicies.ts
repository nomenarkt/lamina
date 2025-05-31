import { useQuery } from '@tanstack/react-query';

export function useUserPolicies(userId?: number, orgUnitId?: number) {
  return useQuery({
    queryKey: ['user-permissions', userId, orgUnitId],
    queryFn: async () => {
      if (!userId || !orgUnitId) return [];
      const res = await fetch(
        `/api/v1/admin/user/${userId}/policies?org_unit_id=${orgUnitId}`
      );
      if (!res.ok) throw new Error('Failed to load user permissions');
      return res.json(); // [["user:42", "/api/flights", "read"], ...]
    },
    enabled: !!userId && !!orgUnitId,
  });
}
