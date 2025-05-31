'use client';

import { useState } from 'react';

type User = { id: number; email: string };

export default function UserSelect({
  value,
  onChange,
  id, // ✅ NEW
}: {
  value: number | undefined;
  onChange: (id: number) => void;
  id?: string; // ✅ NEW
}) {
  const [users] = useState<User[]>([
    { id: 1, email: 'admin@example.com' },
    { id: 2, email: 'auditor@example.com' },
  ]); // Replace with useQuery in real implementation

  return (
    <select
      id={id} // ✅ NEW
      value={value}
      onChange={(e) => onChange(Number(e.target.value))}
      className="border px-3 py-2 rounded"
    >
      <option value="">Select user</option>
      {users.map((user) => (
        <option key={user.id} value={user.id}>
          {user.email}
        </option>
      ))}
    </select>
  );
}
