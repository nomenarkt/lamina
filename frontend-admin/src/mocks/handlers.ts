// src/mocks/handlers.ts
import { http, HttpResponse } from 'msw';

export const handlers = [
  http.get('/api/v1/admin/roles', ({ request }) => {
    const url = new URL(request.url);
    const userId = url.searchParams.get('user_id');
    const orgUnitId = url.searchParams.get('org_unit_id');

    console.log('Mock roles: userId:', userId, 'orgUnitId:', orgUnitId);

    if (userId === '1' && orgUnitId === '100') {
      return HttpResponse.json(
        [{ user_id: 1, org_unit_id: 100, function: 'admin' }],
        { status: 200 }
      );
    }

    return HttpResponse.json([], { status: 200 });
  }),

  http.get('/api/v1/admin/policies', () => {
    console.log('Mock policies hit');
    return HttpResponse.json(
      [
        {
          role: 'admin',
          org_unit_id: 100,
          action: 'read',
          object: '/api/resource',
        },
      ],
      { status: 200 }
    );
  }),
];
