import { TextEncoder, TextDecoder } from 'util';
import { TransformStream } from 'stream/web';

globalThis.TextEncoder = TextEncoder;
globalThis.TextDecoder = TextDecoder;
globalThis.TransformStream = TransformStream;

// Minimal no-op BroadcastChannel polyfill for MSW to avoid runtime error
globalThis.BroadcastChannel = class {
  constructor() {}
  postMessage() {}
  close() {}
  addEventListener() {}
  removeEventListener() {}
};
