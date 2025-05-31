import React from 'react';
import { renderHook, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { useAssignedRoles } from '../useRoles';
import { setupFetchMock } from '@/test/setupFetchMock';

const createTestClient = () =>
  new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
        gcTime: 0,
      },
    },
  });

const wrapper = ({ children }: { children: React.ReactNode }) => (
  <QueryClientProvider client={createTestClient()}>{children}</QueryClientProvider>
);

beforeEach(() => {
  setupFetchMock({
    '/api/v1/admin/roles?user_id=1&org_unit_id=100': {
      body: [{ user_id: 1, org_unit_id: 100, function: 'admin' }],
    },
    '/api/v1/admin/policies': {
      body: [
        {
          role: 'admin',
          org_unit_id: 100,
          object: '/api/resource',
          action: 'read',
        },
      ],
    },
  });
});


describe('useRoles hook', () => {
  test('fetches assigned roles', async () => {
    const { result } = renderHook(() => useAssignedRoles(1, 100), { wrapper });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));

    expect(result.current.data).toEqual([
      {
        user_id: 1,
        org_unit_id: 100,
        function: 'admin',
      },
    ]);
  });
});
