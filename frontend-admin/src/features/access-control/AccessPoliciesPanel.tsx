'use client';

import { useState } from 'react';
import OrgUnitSelect from './components/OrgUnitSelect';
import PolicyTable from './components/PolicyTable';
import { useAddPolicy } from './hooks/usePolicies';

const ACTIONS = ['read', 'write', '*'];
const FUNCTIONS = ['planner', 'auditor', 'admin'];

export default function AccessPoliciesPanel() {
  const [role, setRole] = useState('');
  const [orgUnitId, setOrgUnitId] = useState<number>();
  const [object, setObject] = useState('');
  const [action, setAction] = useState('');
  const [error, setError] = useState('');
  const addPolicy = useAddPolicy();

  const submit = () => {
    setError('');
    if (!role || !orgUnitId || !object || !action) return;
    addPolicy.mutate(
      { role, org_unit_id: orgUnitId, object, action },
      {
        onError: (err) => setError((err as Error).message),
      }
    );
  };

  return (
    <div className="space-y-4">
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div>
          <label htmlFor="role-select" className="block mb-1 font-medium">Role</label>
          <select
            id="role-select"
            value={role}
            onChange={(e) => setRole(e.target.value)}
            className="border px-3 py-2 rounded w-full"
          >
            <option value="">Select role</option>
            {FUNCTIONS.map((fn) => (
              <option key={fn} value={fn}>
                {fn}
              </option>
            ))}
          </select>
        </div>

        <div>
          <label htmlFor="orgunit-select" className="block mb-1 font-medium">Org Unit</label>
          <OrgUnitSelect id="orgunit-select" value={orgUnitId} onChange={setOrgUnitId} />
        </div>

        <div>
          <label htmlFor="object-input" className="block mb-1 font-medium">Object</label>
          <input
            id="object-input"
            value={object}
            onChange={(e) => setObject(e.target.value)}
            className="border px-3 py-2 rounded w-full"
            placeholder="/api/flights"
          />
        </div>

        <div>
          <label htmlFor="action-select" className="block mb-1 font-medium">Action</label>
          <select
            id="action-select"
            value={action}
            onChange={(e) => setAction(e.target.value)}
            className="border px-3 py-2 rounded w-full"
          >
            <option value="">Select action</option>
            {ACTIONS.map((a) => (
              <option key={a} value={a}>
                {a}
              </option>
            ))}
          </select>
        </div>
      </div>

      <button
        onClick={submit}
        className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
      >
        Add Policy
      </button>

      {error && <p className="text-sm text-red-600">{error}</p>}

      <PolicyTable />
    </div>
  );
}
