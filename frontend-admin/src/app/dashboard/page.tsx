'use client';
import { useEffect } from 'react';
import { useRouter } from 'next/navigation';

export default function DashboardIndex() {
  const router = useRouter();

  useEffect(() => {
    const role = localStorage.getItem('user_role');

    if (role === 'admin') router.replace('/dashboard/admin');
    else if (role === 'planner') router.replace('/dashboard/planner');
    else router.replace('/dashboard/crew');
  }, [router]);

  return <p className="text-center mt-8">Redirecting to your dashboard...</p>;
}
