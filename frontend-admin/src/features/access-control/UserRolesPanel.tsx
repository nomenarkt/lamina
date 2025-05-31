'use client';

import { useState } from 'react';
import UserSelect from './components/UserSelect';
import OrgUnitSelect from './components/OrgUnitSelect';
import FunctionMultiSelect from './components/FunctionMultiSelect';
import RoleAssignmentTable from './components/RoleAssignmentTable';
import UserEffectivePermissions from './components/UserEffectivePermissions';
import { useAssignRole } from './hooks/useRoles';

export default function UserRolesPanel() {
  const [userId, setUserId] = useState<number>();
  const [orgUnitId, setOrgUnitId] = useState<number>();
  const [selectedFunctions, setSelectedFunctions] = useState<string[]>([]);
  const [error, setError] = useState('');
  const assignRole = useAssignRole();

  const assign = () => {
    setError('');
    if (!userId || !orgUnitId || selectedFunctions.length === 0) return;
    selectedFunctions.forEach((fn) => {
      assignRole.mutate(
        { user_id: userId, function: fn, org_unit_id: orgUnitId },
        {
          onError: (err) => {
            setError((err as Error).message);
          },
        }
      );
    });
  };

  return (
    <div className="space-y-4">
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div>
          <label htmlFor="user-select" className="block mb-1 font-medium">User</label>
          <UserSelect id="user-select" value={userId} onChange={setUserId} />
        </div>
        <div>
          <label htmlFor="orgunit-select" className="block mb-1 font-medium">Org Unit</label>
          <OrgUnitSelect id="orgunit-select" value={orgUnitId} onChange={setOrgUnitId} />
        </div>
        <div>
          <label htmlFor="function-multiselect" className="block mb-1 font-medium">Functions</label>
          <FunctionMultiSelect selected={selectedFunctions} onChange={setSelectedFunctions} />
        </div>
      </div>

      <button
        onClick={assign}
        className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
      >
        Assign Role(s)
      </button>

      {error && <p className="text-sm text-red-600">{error}</p>}

      <RoleAssignmentTable userId={userId} orgUnitId={orgUnitId} />
      <UserEffectivePermissions userId={userId} orgUnitId={orgUnitId} />
    </div>
  );
}
