import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { SignupForm } from '../SignupForm';


describe('SignupForm', () => {
  it('shows validation error if email is invalid', async () => {
    render(<SignupForm />);

    fireEvent.change(screen.getByPlaceholderText('Email'), {
      target: { value: 'invalid-email' },
    });
    fireEvent.blur(screen.getByPlaceholderText('Email')); // ✨ triggers validation

    fireEvent.change(screen.getByPlaceholderText('Password'), {
      target: { value: 'secure123' },
    });
    fireEvent.change(screen.getByPlaceholderText('Confirm Password'), {
      target: { value: 'secure123' },
    });

    fireEvent.click(screen.getByText('Signup'));

    await waitFor(() => {
      const error = screen.queryByText((content) =>
        content.toLowerCase().includes('invalid') &&
        content.toLowerCase().includes('email')
      );
      expect(error).not.toBeNull();
    });
  });

  it("shows error if passwords don't match", async () => {
    render(<SignupForm />);

    fireEvent.change(screen.getByPlaceholderText('Email'), {
      target: { value: 'user@madagascarairlines.com' },
    });
    fireEvent.change(screen.getByPlaceholderText('Password'), {
      target: { value: 'secure123' },
    });
    fireEvent.change(screen.getByPlaceholderText('Confirm Password'), {
      target: { value: 'mismatch123' },
    });

    fireEvent.click(screen.getByText('Signup'));

    await waitFor(() => {
      expect(
        screen.getByText(/passwords do not match/i)
      ).toBeInTheDocument();
    });
  });
});
