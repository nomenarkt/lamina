import React from 'react';
import { renderHook, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { usePolicies } from '../usePolicies';
import { setupFetchMock } from '@/test/setupFetchMock';

const createTestClient = () =>
  new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
        // 👇 'cacheTime' should be nested inside 'gcTime'
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

describe('usePolicies hook', () => {
  test('fetches policies', async () => {
    const { result } = renderHook(() => usePolicies(), { wrapper });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));

    expect(result.current.data).toEqual([
      {
        role: 'admin',
        org_unit_id: 100,
        object: '/api/resource',
        action: 'read',
      },
    ]);
  });
});
