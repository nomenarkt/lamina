import 'whatwg-fetch';
import '@testing-library/jest-dom';
import { server } from './mocks/server';
import { TextEncoder, TextDecoder } from 'util'; // 👈 Required for MSW 2.x

// ✅ Polyfill for Node.js < 20
if (typeof global.TextEncoder === 'undefined') {
  global.TextEncoder = TextEncoder as typeof global.TextEncoder;
}
if (typeof global.TextDecoder === 'undefined') {
  global.TextDecoder = TextDecoder as typeof global.TextDecoder;
}

// Setup MSW
beforeAll(() => server.listen());
afterEach(() => server.resetHandlers());
afterAll(() => server.close());
