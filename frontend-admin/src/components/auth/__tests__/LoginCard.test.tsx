import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { LoginCard } from '../LoginCard';
import React from 'react';


// Mock the Next.js App Router
jest.mock('next/navigation', () => ({
  useRouter: () => ({
    push: jest.fn(),
  }),
}));

beforeEach(() => {
  global.fetch = jest.fn(() =>
    Promise.resolve({
      ok: false,
      json: async () => ({ error: 'invalid email or password' }),
    })
  );
});

describe('LoginCard', () => {
  it('shows error when fields are empty', async () => {
    render(<LoginCard />);
    fireEvent.click(screen.getByRole('button', { name: /log in/i }));

    await waitFor(() => {
      expect(screen.getByRole('alert')).toHaveTextContent(/please fill in all fields/i);
    });
  });

  it('shows error for invalid credentials', async () => {
    render(<LoginCard />);
    fireEvent.change(screen.getByLabelText(/email/i), { target: { value: 'user@example.com' } });
    fireEvent.change(screen.getByLabelText(/password/i), { target: { value: 'secret' } });
    fireEvent.click(screen.getByRole('button', { name: /log in/i }));

    await waitFor(() => {
      expect(screen.getByRole('alert')).toHaveTextContent(/invalid email or password/i);
    });
  });
});
