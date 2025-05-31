import { render, screen, fireEvent } from '@testing-library/react';
import AccessPoliciesPanel from '../AccessPoliciesPanel';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import React from 'react';


const wrapper = ({ children }: { children: React.ReactNode }) => (
  <QueryClientProvider client={new QueryClient()}>{children}</QueryClientProvider>
);

describe('AccessPoliciesPanel component', () => {
  test('renders form elements', () => {
    render(<AccessPoliciesPanel />, { wrapper });

    expect(screen.getByLabelText(/Role/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/Org Unit/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/Object/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/Action/i)).toBeInTheDocument();
    expect(screen.getByText(/Add Policy/i)).toBeInTheDocument();
  });

  test('add policy button is disabled when inputs are incomplete', () => {
    render(<AccessPoliciesPanel />, { wrapper });

    const addButton = screen.getByText(/Add Policy/i);
    fireEvent.click(addButton);

    // Add assertions based on your implementation
  });

  test('adds policy when inputs are valid', async () => {
    render(<AccessPoliciesPanel />, { wrapper });

    // Simulate user interactions to select role, org unit, object, and action

    const addButton = screen.getByText(/Add Policy/i);
    fireEvent.click(addButton);

    // Add assertions based on your implementation
  });
});
