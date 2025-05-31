'use client';

import { useAssignedRoles, useRemoveRole } from '../hooks/useRoles';
import { getOrgNameById } from '../constants/orgUnits';

export default function RoleAssignmentTable({
  userId,
  orgUnitId,
}: {
  userId?: number;
  orgUnitId?: number;
}) {
  const { data = [], isLoading } = useAssignedRoles(userId, orgUnitId);
  const removeRole = useRemoveRole();

  if (!userId || !orgUnitId) return null;
  if (isLoading) return <p>Loading roles...</p>;

  return (
    <div className="mt-4">
      <p className="text-sm text-gray-600">
        Showing roles for <strong>{getOrgNameById(orgUnitId)}</strong>
      </p>
      <table className="table-auto w-full border mt-2">
        <thead>
          <tr className="bg-gray-100">
            <th className="p-2 text-left">Function</th>
            <th className="p-2 text-center">Actions</th>
          </tr>
        </thead>
        <tbody>
          {data.map((fn: string) => (
            <tr key={fn}>
              <td className="p-2">{fn}</td>
              <td className="p-2 text-center">
                <button
                  className="text-red-600"
                  onClick={() =>
                    removeRole.mutate({ user_id: userId, function: fn, org_unit_id: orgUnitId })
                  }
                >
                  Remove
                </button>
              </td>
            </tr>
          ))}
          {data.length === 0 && (
            <tr>
              <td colSpan={2} className="p-2 text-center text-sm text-gray-500">
                No roles assigned.
              </td>
            </tr>
          )}
        </tbody>
      </table>
    </div>
  );
}
