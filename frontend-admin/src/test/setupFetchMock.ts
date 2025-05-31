type MockResponse = {
  status?: number;
  headers?: Record<string, string>;
  body?: unknown;
};

function jsonResponse(body: unknown, init?: Omit<ResponseInit, 'body'>) {
  return new Response(JSON.stringify(body), {
    status: 200,
    headers: { 'Content-Type': 'application/json', ...(init?.headers || {}) },
    ...init,
  });
}

export function setupFetchMock(handlers: Record<string, MockResponse | (() => MockResponse)>) {
  global.fetch = jest.fn((input: RequestInfo | URL) => {
    const url = typeof input === 'string'
      ? input
      : input instanceof Request
        ? input.url
        : input.toString();

    console.log('FETCHED URL:', url);

    const match = Object.entries(handlers).find(([key]) => url.startsWith(key));

    if (match) {
      const [, handler] = match;
      const response = typeof handler === 'function' ? handler() : handler;
      return Promise.resolve(jsonResponse(response.body, {
        status: response.status,
        headers: response.headers,
      }));
    }

    return Promise.resolve(jsonResponse({}, { status: 200 }));
  }) as typeof fetch;
}
