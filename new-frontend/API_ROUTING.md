# API Routing Configuration

## Overview

The frontend handles API routing differently for development and production to ensure smooth operation in both environments.

## Development Mode (localhost)

**Configuration:**
- `NEXT_PUBLIC_API_URL=http://localhost:8080`
- Direct connection to backend
- **No `/api/` prefix needed**

**How it works:**
```
Frontend (localhost:3000)
    ↓
Direct axios call
    ↓
Backend (localhost:8080/v1/accounts)
```

**Code:**
```typescript
// src/lib/api/client.ts
const baseURL = "http://localhost:8080";  // In dev

// API calls
GET http://localhost:8080/v1/accounts
POST http://localhost:8080/v1/accounts
GET http://localhost:8080/v1/accounts/123/image
```

## Production Mode

**Configuration:**
- `baseURL = "/api"` in axios client
- Next.js rewrites `/api/*` to actual backend

**How it works:**
```
Frontend
    ↓
Relative URL: /api/v1/accounts
    ↓
Next.js rewrite
    ↓
Backend (actual-backend.com/v1/accounts)
```

**Code:**
```typescript
// src/lib/api/client.ts
const baseURL = "/api";  // In production

// API calls become relative
GET /api/v1/accounts
POST /api/v1/accounts

// Next.js rewrites to:
GET https://actual-backend.com/v1/accounts
POST https://actual-backend.com/v1/accounts
```

## Implementation Details

### 1. API Client (`src/lib/api/client.ts`)

```typescript
const isProduction = process.env.NODE_ENV === "production";
const baseURL = isProduction
  ? "/api"  // Production: use Next.js proxy
  : (process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080");  // Dev: direct

export const apiClient = axios.create({
  baseURL,
  headers: { "Content-Type": "application/json" },
});
```

### 2. API Endpoints (no `/api/` in path)

All endpoint paths are defined **without** the `/api/` prefix:

```typescript
// ✅ Correct
GET /v1/accounts
POST /v1/accounts

// ❌ Wrong (old approach)
GET /api/v1/accounts
POST /api/v1/accounts
```

### 3. Image URLs Utility (`src/lib/utils.ts`)

For `<img>` tags that don't go through axios:

```typescript
export function getApiUrl(path: string): string {
  const isProduction = process.env.NODE_ENV === "production";
  const baseUrl = isProduction
    ? "/api"
    : (process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080");
  return `${baseUrl}${path}`;
}

// Usage
<img src={getApiUrl(`/v1/accounts/${id}/image`)} />
```

### 4. Next.js Rewrites (`next.config.js`)

```javascript
async rewrites() {
  if (process.env.NODE_ENV === 'production' && process.env.NEXT_PUBLIC_API_URL) {
    return [
      {
        source: '/api/:path*',
        destination: `${process.env.NEXT_PUBLIC_API_URL}/:path*`,
      },
    ];
  }
  return [];
}
```

## Environment Variables

### Development (`.env.local`)
```env
NEXT_PUBLIC_API_URL=http://localhost:8080
```

### Production (`.env.production`)
```env
NEXT_PUBLIC_API_URL=https://api.geekbudget.com
```

## Benefits

### Development
✅ Direct connection to backend (faster, simpler debugging)
✅ No proxy overhead
✅ Works with backend on different port
✅ No CORS issues (backend must allow localhost:3000)

### Production
✅ No CORS issues (same-origin via Next.js proxy)
✅ Backend URL hidden from client
✅ Can change backend URL without redeploying frontend
✅ Works with API gateways and load balancers

## Testing

### Development Server
```bash
# Terminal 1: Start backend
cd backend && make run-backend
# Backend runs on http://localhost:8080

# Terminal 2: Start frontend
cd new-frontend && npm run dev
# Frontend runs on http://localhost:3000

# Test API call
curl http://localhost:3000
# Frontend makes request to: http://localhost:8080/v1/accounts
```

### Production Build
```bash
npm run build
npm start

# Test API call
# Frontend makes request to: /api/v1/accounts
# Next.js rewrites to: ${NEXT_PUBLIC_API_URL}/v1/accounts
```

## Troubleshooting

### CORS errors in development
**Problem:** Backend rejects requests from `localhost:3000`
**Solution:** Configure backend to allow `http://localhost:3000` in CORS settings

### 404 errors in development
**Problem:** API calls return 404
**Solution:** Check that backend is running on port 8080 and paths don't include `/api/` prefix

### Connection refused
**Problem:** `ECONNREFUSED localhost:8080`
**Solution:** Start the backend server first

### Image loading fails
**Problem:** Account images don't load
**Solution:** Ensure `getApiUrl()` utility is used for image URLs, not hardcoded paths
