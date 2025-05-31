/**
 * @jest-environment jsdom
 */

import { render, screen } from '@testing-library/react';
import ConfirmErrorPage from '@/app/confirm-error/page';
import { useSearchParams } from 'next/navigation';
import React from 'react';


jest.mock('next/navigation', () => ({
  useSearchParams: jest.fn(),
}));

describe('ConfirmErrorPage', () => {
  it('displays expired message and Signup link', () => {
    (useSearchParams as jest.Mock).mockReturnValue({
      get: () => 'expired',
    });

    render(<ConfirmErrorPage />);
    expect(screen.getByText(/your confirmation link has expired/i)).toBeInTheDocument();
    expect(screen.getByRole('link', { name: /return to signup/i })).toHaveAttribute('href', '/signup');
  });

  it('displays already-confirmed message and Login link', () => {
    (useSearchParams as jest.Mock).mockReturnValue({
      get: () => 'already-confirmed',
    });

    render(<ConfirmErrorPage />);
    expect(screen.getByText(/already been confirmed/i)).toBeInTheDocument();
    expect(screen.getByRole('link', { name: /return to login/i })).toHaveAttribute('href', '/login');
  });

  it('displays invalid message and Signup link', () => {
    (useSearchParams as jest.Mock).mockReturnValue({
      get: () => 'invalid',
    });

    render(<ConfirmErrorPage />);
    expect(screen.getByText(/invalid or malformed confirmation link/i)).toBeInTheDocument();
    expect(screen.getByRole('link', { name: /return to signup/i })).toHaveAttribute('href', '/signup');
  });

  it('displays fallback error message and Signup link when no reason provided', () => {
    (useSearchParams as jest.Mock).mockReturnValue({
      get: () => null,
    });

    render(<ConfirmErrorPage />);
    expect(screen.getByText(/something went wrong/i)).toBeInTheDocument();
    expect(screen.getByRole('link', { name: /return to signup/i })).toHaveAttribute('href', '/signup');
  });
});
