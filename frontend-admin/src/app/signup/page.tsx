import { SignupForm } from '@/features/auth/components/SignupForm';

export default function SignupPage() {
  return (
    <div className="max-w-md mx-auto p-6">
      <h1 className="text-2xl font-bold mb-4">Signup</h1>
      <SignupForm />
    </div>
  );
}
