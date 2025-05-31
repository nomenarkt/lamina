// frontend-admin/jest.config.cjs
module.exports = {
  transform: {
    '^.+\\.(ts|tsx)$': ['babel-jest', { configFile: './jest.babel.config.js' }],
  },
  testEnvironment: 'jsdom',
  testEnvironmentOptions: {
    customExportConditions: [],
  },
  setupFiles: ['<rootDir>/jest.polyfill.ts'], // 👈 Runs before any test files
  setupFilesAfterEnv: ['<rootDir>/src/setupTests.ts'],
  moduleNameMapper: {
    '^@/(.*)$': '<rootDir>/src/$1',
  },
  testMatch: ['<rootDir>/src/**/*.test.(ts|tsx)'],
  testEnvironment: 'jest-fixed-jsdom',
};
