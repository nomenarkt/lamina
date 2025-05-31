// jest.setup.ts
import '@testing-library/jest-dom';
import 'whatwg-fetch';
import { server } from './src/mocks/server';
import { TextEncoder, TextDecoder } from 'util';
import { BroadcastChannel } from 'worker_threads'; // Polyfill for BroadcastChannel

// Ensure global availability for MSW compatibility
Object.assign(globalThis, {
  TextEncoder,
  TextDecoder,
  BroadcastChannel,
});

// MSW server lifecycle hooks
beforeAll(() => server.listen());
afterEach(() => server.resetHandlers());
afterAll(() => server.close());
