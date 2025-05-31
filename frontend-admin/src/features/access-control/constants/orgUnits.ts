export const ORG_UNITS = [
  { id: 47, name: 'Exploitation Aérienne' },
  { id: 48, name: 'Finance' },
];

export function getOrgNameById(id?: number | string): string {
  if (!id) return '';
  const match = ORG_UNITS.find((o) => Number(o.id) === Number(id));
  return match?.name ?? `orgunit:${id}`;
}
