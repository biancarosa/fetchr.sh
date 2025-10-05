# Playwright E2E Tests

This directory contains end-to-end tests for the Netkit dashboard using Playwright.

## Prerequisites

Before running the tests, you need to have:

1. **Netkit backend running** on port 8080 (proxy) and 8081 (admin)
2. **Dashboard dev server running** on port 3000 (or built dashboard)

## Running Tests

### Option 1: Automatic (Recommended for Local Development)

The Playwright config will automatically start the dev server for you:

```bash
# From the dashboard directory
npm run test:e2e
```

**Note:** You still need to start the netkit backend separately:

```bash
# From the root directory, in a separate terminal
make serve-dev
```

### Option 2: Manual Setup

1. Start the netkit backend:
   ```bash
   # From root directory
   make serve-dev
   ```

2. In another terminal, start the dashboard:
   ```bash
   # From root directory
   make dashboard-dev
   ```

3. In a third terminal, run the tests:
   ```bash
   # From dashboard directory
   npm run test:e2e
   ```

### Option 3: Using Make (Full Stack)

From the root directory:

```bash
# Start both backend and dashboard
make dev
```

Then in another terminal:

```bash
# Run E2E tests
make test-e2e
```

## Test Commands

- `npm run test:e2e` - Run tests headless (CI mode)
- `npm run test:e2e:ui` - Run tests with Playwright UI (interactive)
- `npm run test:e2e:headed` - Run tests with visible browser
- `npm run test:e2e:report` - View last test report

Or from root directory with Make:
- `make test-e2e` - Run tests headless
- `make test-e2e-ui` - Run with UI
- `make test-e2e-headed` - Run with visible browser

## CI/CD

In GitHub Actions, the workflow:
1. Builds the netkit backend
2. Builds the dashboard
3. Starts the netkit server on ports 8080/8081/3000
4. Runs the Playwright tests
5. Uploads test reports as artifacts

The tests will run automatically on:
- Push to main branch
- Pull requests to main branch

## Writing Tests

Tests are located in `dashboard/e2e/*.spec.ts`. The tests:

- Use Playwright's page object model
- Target the dashboard at `http://localhost:3000`
- Assume the proxy is running at `http://localhost:8080`
- Assume the admin API is at `http://localhost:8081`

Example test structure:

```typescript
import { test, expect } from '@playwright/test';

test.describe('Feature Name', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
  });

  test('should do something', async ({ page }) => {
    // Your test code
  });
});
```

## Debugging Failed Tests

If tests fail:

1. Check the HTML report: `npm run test:e2e:report`
2. Run with UI mode: `npm run test:e2e:ui`
3. Check screenshots in `test-results/` directory
4. Verify backend is running: `curl http://localhost:8081/healthz`
5. Verify dashboard is running: `curl http://localhost:3000`

## Common Issues

### "Proxy Offline" error
- The dashboard shows "Proxy Offline" when it can't reach the backend
- Ensure `make serve-dev` is running
- Check ports 8080 and 8081 are not in use by other processes

### Tests timeout
- Increase timeout in `playwright.config.ts` if needed
- Check network connectivity
- Verify the backend responds quickly: `time curl http://localhost:8081/healthz`

### Element not found errors
- The dashboard UI may have changed
- Update selectors in test files
- Use Playwright Inspector: `npx playwright test --debug`
