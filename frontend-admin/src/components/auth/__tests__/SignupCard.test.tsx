import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { SignupCard } from '../SignupCard';
import React from 'react';


// Mock the Next.js App Router
jest.mock('next/navigation', () => ({
  useRouter: () => ({
    push: jest.fn(),
  }),
}));

describe('SignupCard', () => {
  it('shows error when fields are empty', async () => {
    render(<SignupCard />);
    fireEvent.click(screen.getByRole('button', { name: /sign up/i }));

    await waitFor(() => {
      expect(screen.getByRole('alert')).toHaveTextContent(/please fill in all fields/i);
    });
  });

  it('shows error when passwords do not match', async () => {
    render(<SignupCard />);
    fireEvent.change(screen.getByLabelText(/^email$/i), { target: { value: 'test@madagascarairlines.com' } });
    fireEvent.change(screen.getByLabelText(/^password$/i), { target: { value: 'abc123' } });
    fireEvent.change(screen.getByLabelText(/confirm password/i), { target: { value: 'xyz789' } });

    fireEvent.click(screen.getByRole('button', { name: /sign up/i }));

    await waitFor(() => {
      expect(screen.getByRole('alert')).toHaveTextContent(/passwords do not match/i);
    });
  });
});
