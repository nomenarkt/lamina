import { InviteUserForm } from '@/features/users/components/InviteUserForm';

export default function InvitePage() {
  return (
    <main className="p-8 max-w-xl mx-auto">
      <h1 className="text-2xl font-bold mb-4">Invite New User</h1>
      <InviteUserForm />
    </main>
  );
}
