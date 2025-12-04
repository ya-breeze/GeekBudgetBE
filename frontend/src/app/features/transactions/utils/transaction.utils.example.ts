/**
 * Example usage of TransactionUtils
 * 
 * This file demonstrates how to use the TransactionUtils class
 * to extract input and output accounts from transactions.
 */

import { Transaction } from '../../../core/api/models/transaction';
import { TransactionUtils } from './transaction.utils';

// Example 1: Simple transfer between two accounts
export function exampleSimpleTransfer() {
  const transaction: Transaction = {
    id: '1',
    date: '2024-01-15',
    description: 'Transfer to savings',
    movements: [
      { accountId: 'checking-account-id', amount: -500, currencyId: 'usd' },
      { accountId: 'savings-account-id', amount: 500, currencyId: 'usd' },
    ],
  };

  const inputAccount = TransactionUtils.getInputAccountId(transaction);
  const outputAccount = TransactionUtils.getOutputAccountId(transaction);

  console.log('Transfer from:', inputAccount); // 'checking-account-id'
  console.log('Transfer to:', outputAccount);   // 'savings-account-id'
}

// Example 2: Expense transaction
export function exampleExpense() {
  const transaction: Transaction = {
    id: '2',
    date: '2024-01-16',
    description: 'Grocery shopping',
    movements: [
      { accountId: 'checking-account-id', amount: -150.50, currencyId: 'usd' },
      { accountId: 'groceries-expense-id', amount: 150.50, currencyId: 'usd' },
    ],
  };

  const sourceAccount = TransactionUtils.getInputAccountId(transaction);
  const expenseCategory = TransactionUtils.getOutputAccountId(transaction);

  console.log('Paid from:', sourceAccount);      // 'checking-account-id'
  console.log('Expense category:', expenseCategory); // 'groceries-expense-id'
}

// Example 3: Split payment from multiple sources
export function exampleSplitPayment() {
  const transaction: Transaction = {
    id: '3',
    date: '2024-01-17',
    description: 'Restaurant bill split between accounts',
    movements: [
      { accountId: 'checking-account-id', amount: -60, currencyId: 'usd' },
      { accountId: 'credit-card-id', amount: -40, currencyId: 'usd' },
      { accountId: 'restaurant-expense-id', amount: 100, currencyId: 'usd' },
    ],
  };

  const sourceAccounts = TransactionUtils.getInputAccountIds(transaction);
  const expenseCategory = TransactionUtils.getOutputAccountId(transaction);

  console.log('Paid from accounts:', sourceAccounts); // ['checking-account-id', 'credit-card-id']
  console.log('Expense category:', expenseCategory);  // 'restaurant-expense-id'
}

// Example 4: Split expense to multiple categories
export function exampleSplitExpense() {
  const transaction: Transaction = {
    id: '4',
    date: '2024-01-18',
    description: 'Shopping - groceries and household items',
    movements: [
      { accountId: 'checking-account-id', amount: -200, currencyId: 'usd' },
      { accountId: 'groceries-expense-id', amount: 120, currencyId: 'usd' },
      { accountId: 'household-expense-id', amount: 80, currencyId: 'usd' },
    ],
  };

  const sourceAccount = TransactionUtils.getInputAccountId(transaction);
  const expenseCategories = TransactionUtils.getOutputAccountIds(transaction);

  console.log('Paid from:', sourceAccount);           // 'checking-account-id'
  console.log('Expense categories:', expenseCategories); // ['groceries-expense-id', 'household-expense-id']
}

// Example 5: Income transaction
export function exampleIncome() {
  const transaction: Transaction = {
    id: '5',
    date: '2024-01-20',
    description: 'Monthly salary',
    movements: [
      { accountId: 'salary-income-id', amount: -3000, currencyId: 'usd' },
      { accountId: 'checking-account-id', amount: 3000, currencyId: 'usd' },
    ],
  };

  const incomeSource = TransactionUtils.getInputAccountId(transaction);
  const destinationAccount = TransactionUtils.getOutputAccountId(transaction);

  console.log('Income from:', incomeSource);        // 'salary-income-id'
  console.log('Deposited to:', destinationAccount); // 'checking-account-id'
}

// Example 6: Using with account map to get account names
export function exampleWithAccountMap(transaction: Transaction, accountMap: Map<string, { name: string }>) {
  const inputAccountId = TransactionUtils.getInputAccountId(transaction);
  const outputAccountId = TransactionUtils.getOutputAccountId(transaction);

  const inputAccountName = inputAccountId ? accountMap.get(inputAccountId)?.name : 'Unknown';
  const outputAccountName = outputAccountId ? accountMap.get(outputAccountId)?.name : 'Unknown';

  return {
    from: inputAccountName,
    to: outputAccountName,
  };
}

// Example 7: Formatting transaction for display
export function formatTransactionSummary(
  transaction: Transaction,
  accountMap: Map<string, { name: string }>
): string {
  const inputIds = TransactionUtils.getInputAccountIds(transaction);
  const outputIds = TransactionUtils.getOutputAccountIds(transaction);

  if (inputIds.length === 0 || outputIds.length === 0) {
    return 'Invalid transaction';
  }

  const inputNames = inputIds.map(id => accountMap.get(id)?.name || id);
  const outputNames = outputIds.map(id => accountMap.get(id)?.name || id);

  if (inputIds.length === 1 && outputIds.length === 1) {
    return `${inputNames[0]} → ${outputNames[0]}`;
  } else if (inputIds.length > 1 && outputIds.length === 1) {
    return `${inputNames.join(' + ')} → ${outputNames[0]}`;
  } else if (inputIds.length === 1 && outputIds.length > 1) {
    return `${inputNames[0]} → ${outputNames.join(' + ')}`;
  } else {
    return `${inputNames.join(' + ')} → ${outputNames.join(' + ')}`;
  }
}

