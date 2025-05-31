'use client';

const FUNCTIONS = ['planner', 'auditor', 'admin'];

export default function FunctionMultiSelect({
  selected,
  onChange,
}: {
  selected: string[];
  onChange: (newList: string[]) => void;
}) {
  const toggle = (role: string) => {
    if (selected.includes(role)) {
      onChange(selected.filter((r) => r !== role));
    } else {
      onChange([...selected, role]);
    }
  };

  return (
    <div className="flex gap-4 flex-wrap">
      {FUNCTIONS.map((role) => (
        <label key={role} className="flex items-center gap-2">
          <input
            type="checkbox"
            checked={selected.includes(role)}
            onChange={() => toggle(role)}
          />
          {role}
        </label>
      ))}
    </div>
  );
}
