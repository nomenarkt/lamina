import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import CheckEmail from '../page';

// Mock useRouter and useSearchParams
jest.mock('next/navigation', () => ({
  useRouter: () => ({ push: jest.fn() }),
  useSearchParams: () => ({
    get: () => 'crew@madagascarairlines.com',
  }),
}));

describe('CheckEmail page', () => {
  beforeEach(() => {
    jest.resetAllMocks();
  });

  it('renders the confirmation message with email', () => {
    render(<CheckEmail />);
    expect(screen.getByText(/We.+sent.+confirmation/i)).toBeInTheDocument();
    expect(screen.getByText(/crew@madagascarairlines.com/i)).toBeInTheDocument();
  });

  it('resend button triggers fetch and success message', async () => {
    global.fetch = jest.fn(() =>
      Promise.resolve({ ok: true, json: () => Promise.resolve({ message: 'ok' }) })
    ) as jest.Mock;

    render(<CheckEmail />);
    fireEvent.click(screen.getByText(/Resend confirmation/i));

    await waitFor(() =>
      expect(screen.getByText(/new confirmation email/i)).toBeInTheDocument()
    );
  });
});
