# Transaction Utilities

This module provides utility functions for working with Transaction objects.

## TransactionUtils

A utility class that provides methods to extract input and output accounts from transactions based on movement amounts.

### Concepts

- **Input Account**: An account that is the source of money (has negative amount in movements)
- **Output Account**: An account that is the destination of money (has positive amount in movements)

### Methods

#### `getInputAccountId(transaction: Transaction): string | undefined`

Returns the first account ID that is the source of money (negative amount).

**Example:**
```typescript
import { TransactionUtils } from './features/transactions/utils';

const transaction: Transaction = {
  id: '1',
  date: '2024-01-01',
  movements: [
    { accountId: 'checking-account', amount: -100, currencyId: 'usd' },
    { accountId: 'groceries', amount: 100, currencyId: 'usd' },
  ],
};

const inputAccountId = TransactionUtils.getInputAccountId(transaction);
// Returns: 'checking-account'
```

#### `getInputAccountIds(transaction: Transaction): string[]`

Returns all account IDs that are sources of money (negative amounts).

**Example:**
```typescript
const transaction: Transaction = {
  id: '1',
  date: '2024-01-01',
  movements: [
    { accountId: 'checking-account', amount: -50, currencyId: 'usd' },
    { accountId: 'savings-account', amount: -50, currencyId: 'usd' },
    { accountId: 'groceries', amount: 100, currencyId: 'usd' },
  ],
};

const inputAccountIds = TransactionUtils.getInputAccountIds(transaction);
// Returns: ['checking-account', 'savings-account']
```

#### `getOutputAccountId(transaction: Transaction): string | undefined`

Returns the first account ID that is the destination of money (positive amount).

**Example:**
```typescript
const transaction: Transaction = {
  id: '1',
  date: '2024-01-01',
  movements: [
    { accountId: 'checking-account', amount: -100, currencyId: 'usd' },
    { accountId: 'groceries', amount: 100, currencyId: 'usd' },
  ],
};

const outputAccountId = TransactionUtils.getOutputAccountId(transaction);
// Returns: 'groceries'
```

#### `getOutputAccountIds(transaction: Transaction): string[]`

Returns all account IDs that are destinations of money (positive amounts).

**Example:**
```typescript
const transaction: Transaction = {
  id: '1',
  date: '2024-01-01',
  movements: [
    { accountId: 'checking-account', amount: -100, currencyId: 'usd' },
    { accountId: 'groceries', amount: 50, currencyId: 'usd' },
    { accountId: 'transportation', amount: 50, currencyId: 'usd' },
  ],
};

const outputAccountIds = TransactionUtils.getOutputAccountIds(transaction);
// Returns: ['groceries', 'transportation']
```

## Use Cases

### Simple Transfer
```typescript
// Transfer from checking to savings
const transaction: Transaction = {
  id: '1',
  date: '2024-01-01',
  movements: [
    { accountId: 'checking', amount: -500, currencyId: 'usd' },
    { accountId: 'savings', amount: 500, currencyId: 'usd' },
  ],
};

const from = TransactionUtils.getInputAccountId(transaction);  // 'checking'
const to = TransactionUtils.getOutputAccountId(transaction);   // 'savings'
```

### Split Payment
```typescript
// Payment split between two accounts
const transaction: Transaction = {
  id: '2',
  date: '2024-01-01',
  movements: [
    { accountId: 'checking', amount: -60, currencyId: 'usd' },
    { accountId: 'credit-card', amount: -40, currencyId: 'usd' },
    { accountId: 'restaurant', amount: 100, currencyId: 'usd' },
  ],
};

const sources = TransactionUtils.getInputAccountIds(transaction);
// Returns: ['checking', 'credit-card']

const destination = TransactionUtils.getOutputAccountId(transaction);
// Returns: 'restaurant'
```

## Notes

- Methods return `undefined` or empty arrays when no matching movements are found
- Movements without an `accountId` are filtered out
- The order of accounts in the arrays matches the order in the movements array

