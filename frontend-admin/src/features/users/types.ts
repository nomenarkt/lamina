export type UserRole = 'admin' | 'viewer' | 'external';

export interface InviteUserPayload {
  email: string;
  role: UserRole;
  company?: string;
  accessDuration?: { from: string; to: string }; // ISO 8601
}
