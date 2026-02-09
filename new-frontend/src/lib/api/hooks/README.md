# Custom API Hooks

This directory contains **non-generated** custom hooks that wrap the generated API client code.

## Why?

The generated API code (in `../generated/`) is auto-generated from the OpenAPI spec and will be overwritten on regeneration. However, it may have issues:

1. **Response type mismatches** - The generator assumes all responses are wrapped in `{data: ...}`, but many endpoints return arrays directly
2. **Inconsistent patterns** - Some endpoints need custom handling
3. **Better DX** - We can provide cleaner APIs with better type inference

## Pattern

For each resource (accounts, transactions, etc.), create a custom hook file:

```typescript
// use-accounts.ts
import { useQuery, useMutation } from "@tanstack/react-query";
import { apiClient } from "../client";
import type { Account } from "../models";

// Define query keys
export const accountKeys = {
  all: ["accounts"] as const,
  lists: () => [...accountKeys.all, "list"] as const,
  // ... more keys
};

// Define API functions with correct types
const getAccounts = async (): Promise<Account[]> => {
  const { data } = await apiClient.get<Account[]>("/v1/accounts");
  return data; // Returns array directly, not {data: array}
};

// Export custom hooks
export function useAccounts() {
  return useQuery({
    queryKey: accountKeys.lists(),
    queryFn: getAccounts,
  });
}
```

## Usage

In your components, **always import from this directory**, not from `../generated/`:

```typescript
// ✅ GOOD - uses custom hooks
import { useAccounts, useCreateAccount } from "@/lib/api/hooks/use-accounts";

// ❌ BAD - uses generated hooks (will break)
import { useGetAccounts } from "@/lib/api/generated/accounts/accounts";
```

## Benefits

- ✅ **Safe regeneration** - Generated code can be regenerated without breaking your app
- ✅ **Correct types** - Match actual API responses
- ✅ **Auto-invalidation** - Mutations automatically invalidate related queries
- ✅ **Cleaner API** - Simpler function signatures
- ✅ **Better DX** - TypeScript inference works better

## When to Create Custom Hooks

Create custom hooks for:
1. Resources used in multiple components (accounts, transactions, etc.)
2. Endpoints with response type mismatches
3. Complex mutation patterns needing query invalidation

For one-off API calls, you can still use the generated client directly.
