'use client';

import { usePolicies, useDeletePolicy } from '../hooks/usePolicies';
import { getOrgNameById } from '../constants/orgUnits';

export default function PolicyTable() {
  const { data = [], isLoading } = usePolicies();
  const deletePolicy = useDeletePolicy();

  if (isLoading) return <p>Loading policies...</p>;

  return (
    <table className="table-auto w-full border mt-4">
      <thead>
        <tr className="bg-gray-100">
          <th className="p-2 text-left">Role</th>
          <th className="p-2 text-left">Org Unit</th>
          <th className="p-2 text-left">Object</th>
          <th className="p-2 text-left">Action</th>
          <th className="p-2 text-center">Actions</th>
        </tr>
      </thead>
      <tbody>
        {data.map(([role, domain, object, action]: string[]) => {
          const orgUnitId = parseInt(domain.replace('orgunit:', ''));
          return (
            <tr key={`${role}-${object}-${action}`}>
              <td className="p-2">{role}</td>
              <td className="p-2">{getOrgNameById(orgUnitId)}</td>
              <td className="p-2">{object}</td>
              <td className="p-2">{action}</td>
              <td className="p-2 text-center">
                <button
                  className="text-red-600"
                  onClick={() =>
                    deletePolicy.mutate({ role, org_unit_id: orgUnitId, object, action })
                  }
                >
                  Delete
                </button>
              </td>
            </tr>
          );
        })}
        {data.length === 0 && (
          <tr>
            <td colSpan={5} className="p-2 text-center text-sm text-gray-500">
              No policies found.
            </td>
          </tr>
        )}
      </tbody>
    </table>
  );
}
