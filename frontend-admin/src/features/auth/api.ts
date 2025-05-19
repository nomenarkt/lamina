import type { LoginPayload, SignupPayload, AuthResponse } from './types';

const API = process.env.NEXT_PUBLIC_API_BASE || 'http://localhost:8080';

export async function login(payload: LoginPayload): Promise<AuthResponse> {
  const res = await fetch(`${API}/api/v1/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  });

  if (!res.ok) {
    throw new Error('Login failed');
  }

  return res.json();
}

export async function signup(payload: SignupPayload): Promise<{ message: string }> {
  const res = await fetch(`${API}/api/v1/auth/signup`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  });

  if (!res.ok) {
    const errorData = await res.json();
    throw new Error(errorData?.message || 'Signup failed');
  }

  return res.json();
}
