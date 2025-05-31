import { render, screen, fireEvent } from '@testing-library/react';
import UserRolesPanel from '../UserRolesPanel';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import React from 'react';


const wrapper = ({ children }: { children: React.ReactNode }) => (
  <QueryClientProvider client={new QueryClient()}>{children}</QueryClientProvider>
);

describe('UserRolesPanel component', () => {
  test('renders form elements', () => {
    render(<UserRolesPanel />, { wrapper });

    expect(screen.getByLabelText(/User/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/Org Unit/i)).toBeInTheDocument();
    expect(screen.getByText(/Functions/i)).toBeInTheDocument();
    expect(screen.getByText(/Assign Role\(s\)/i)).toBeInTheDocument();
  });

  test('assign button is disabled when inputs are incomplete', () => {
    render(<UserRolesPanel />, { wrapper });

    const assignButton = screen.getByText(/Assign Role\(s\)/i);
    fireEvent.click(assignButton);

    // Add assertions based on your implementation
  });

  test('assigns roles when inputs are valid', async () => {
    render(<UserRolesPanel />, { wrapper });

    // Simulate user interactions to select user, org unit, and functions

    const assignButton = screen.getByText(/Assign Role\(s\)/i);
    fireEvent.click(assignButton);

    // Add assertions based on your implementation
  });
});
