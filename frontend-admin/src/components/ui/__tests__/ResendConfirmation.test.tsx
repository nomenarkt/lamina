import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import ResendConfirmation from '../ResendConfirmation';
import React from 'react';


const mockPush = jest.fn();

jest.mock('next/navigation', () => ({
  useRouter: () => ({ push: mockPush }),
}));

beforeEach(() => {
  jest.clearAllMocks();
});

describe('ResendConfirmation', () => {
  const email = 'test@madagascarairlines.com';

  it('renders with default text', () => {
    render(<ResendConfirmation email={email} />);
    expect(screen.getByText('Didn’t get it?')).toBeInTheDocument();
    expect(screen.getByText('Resend confirmation')).toBeInTheDocument();
  });

  it('handles success response', async () => {
    global.fetch = jest.fn(() =>
      Promise.resolve({ ok: true } as Response)
    ) as jest.Mock;

    render(<ResendConfirmation email={email} />);
    fireEvent.click(screen.getByText('Resend confirmation'));

    await waitFor(() =>
      expect(screen.getByText("✅ We've sent you a new confirmation email.")).toBeInTheDocument()
    );
  });

  it('redirects to login if already confirmed (400)', async () => {
    global.fetch = jest.fn(() =>
      Promise.resolve({ ok: false, status: 400 } as Response)
    ) as jest.Mock;

    render(<ResendConfirmation email={email} />);
    fireEvent.click(screen.getByText('Resend confirmation'));

    await waitFor(() => {
      expect(mockPush).toHaveBeenCalledWith('/login');
    });
  });

  it('shows internal-only error (403)', async () => {
    global.fetch = jest.fn(() =>
      Promise.resolve({ ok: false, status: 403 } as Response)
    ) as jest.Mock;

    render(<ResendConfirmation email={email} />);
    fireEvent.click(screen.getByText('Resend confirmation'));

    await waitFor(() =>
      expect(screen.getByText(/only for Madagascar Airlines crew/)).toBeInTheDocument()
    );
  });

  it('redirects to signup if not found (404)', async () => {
    global.fetch = jest.fn(() =>
      Promise.resolve({ ok: false, status: 404 } as Response)
    ) as jest.Mock;

    render(<ResendConfirmation email={email} />);
    fireEvent.click(screen.getByText('Resend confirmation'));

    await waitFor(() => {
      expect(mockPush).toHaveBeenCalledWith('/signup?resend=failed');
    });
  });

  it('handles unknown error', async () => {
    global.fetch = jest.fn(() =>
      Promise.resolve({ ok: false, status: 500 } as Response)
    ) as jest.Mock;

    render(<ResendConfirmation email={email} />);
    fireEvent.click(screen.getByText('Resend confirmation'));

    await waitFor(() =>
      expect(screen.getByText('❌ Could not resend. Please try again later.')).toBeInTheDocument()
    );
  });

  it('handles network error', async () => {
    global.fetch = jest.fn(() => Promise.reject()) as jest.Mock;

    render(<ResendConfirmation email={email} />);
    fireEvent.click(screen.getByText('Resend confirmation'));

    await waitFor(() =>
      expect(screen.getByText('❌ Something went wrong. Please try again later.')).toBeInTheDocument()
    );
  });
});
