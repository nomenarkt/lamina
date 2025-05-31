'use client';

import { useUserPolicies } from '../hooks/useUserPolicies';

export default function UserEffectivePermissions({
  userId,
  orgUnitId,
}: {
  userId?: number;
  orgUnitId?: number;
}) {
  const { data = [], isLoading } = useUserPolicies(userId, orgUnitId);

  if (!userId || !orgUnitId) return null;
  if (isLoading) return <p>Loading effective permissions...</p>;

  return (
    <div className="mt-6 border-t pt-4">
      <h3 className="font-medium mb-2 text-sm text-gray-700">Effective Permissions</h3>
      {data.length === 0 ? (
        <p className="text-sm text-gray-500">No permissions granted.</p>
      ) : (
        <ul className="text-sm list-disc pl-5 space-y-1">
          {data.map(([, obj, act]: [string, string, string], i: number) => (
            <li key={i}>
              {act.toUpperCase()} <code>{obj}</code>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
