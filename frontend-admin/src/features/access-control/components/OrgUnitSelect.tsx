'use client';

type OrgUnit = { id: number; name: string };

const ORG_UNITS: OrgUnit[] = [
  { id: 47, name: 'Exploitation Aérienne' },
  { id: 48, name: 'Finance' },
];

export default function OrgUnitSelect({
  value,
  onChange,
  id, // ✅ NEW
}: {
  value: number | undefined;
  onChange: (id: number) => void;
  id?: string; // ✅ NEW
}) {
  return (
    <select
      id={id} // ✅ NEW
      value={value}
      onChange={(e) => onChange(Number(e.target.value))}
      className="border px-3 py-2 rounded"
    >
      <option value="">Select org unit</option>
      {ORG_UNITS.map((org) => (
        <option key={org.id} value={org.id}>
          {org.name}
        </option>
      ))}
    </select>
  );
}
