module.exports = {
  transform: {
    '^.+\\.(ts|tsx)$': ['babel-jest', { configFile: './jest.babel.config.js' }],
  },
  testEnvironment: 'jsdom',
  setupFilesAfterEnv: ['<rootDir>/jest.setup.ts'],
  moduleNameMapper: {
    '^@/(.*)$': '<rootDir>/src/$1',
  },
  testMatch: ['<rootDir>/src/**/*.test.(ts|tsx)'],
};
