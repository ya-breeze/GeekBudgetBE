import { TransactionUtils } from './transaction.utils';
import { Transaction } from '../../../core/api/models/transaction';

describe('TransactionUtils', () => {
  describe('getInputAccountId', () => {
    it('should return the account ID with negative amount', () => {
      const transaction: Transaction = {
        id: '1',
        date: '2024-01-01',
        movements: [
          { accountId: 'acc1', amount: -100, currencyId: 'usd' },
          { accountId: 'acc2', amount: 100, currencyId: 'usd' },
        ],
      };

      const result = TransactionUtils.getInputAccountId(transaction);
      expect(result).toBe('acc1');
    });

    it('should return undefined when no movements exist', () => {
      const transaction: Transaction = {
        id: '1',
        date: '2024-01-01',
        movements: [],
      };

      const result = TransactionUtils.getInputAccountId(transaction);
      expect(result).toBeUndefined();
    });

    it('should return undefined when no negative amounts exist', () => {
      const transaction: Transaction = {
        id: '1',
        date: '2024-01-01',
        movements: [
          { accountId: 'acc1', amount: 100, currencyId: 'usd' },
          { accountId: 'acc2', amount: 50, currencyId: 'usd' },
        ],
      };

      const result = TransactionUtils.getInputAccountId(transaction);
      expect(result).toBeUndefined();
    });

    it('should return the first input account when multiple exist', () => {
      const transaction: Transaction = {
        id: '1',
        date: '2024-01-01',
        movements: [
          { accountId: 'acc1', amount: -50, currencyId: 'usd' },
          { accountId: 'acc2', amount: -50, currencyId: 'usd' },
          { accountId: 'acc3', amount: 100, currencyId: 'usd' },
        ],
      };

      const result = TransactionUtils.getInputAccountId(transaction);
      expect(result).toBe('acc1');
    });
  });

  describe('getInputAccountIds', () => {
    it('should return all account IDs with negative amounts', () => {
      const transaction: Transaction = {
        id: '1',
        date: '2024-01-01',
        movements: [
          { accountId: 'acc1', amount: -50, currencyId: 'usd' },
          { accountId: 'acc2', amount: -50, currencyId: 'usd' },
          { accountId: 'acc3', amount: 100, currencyId: 'usd' },
        ],
      };

      const result = TransactionUtils.getInputAccountIds(transaction);
      expect(result).toEqual(['acc1', 'acc2']);
    });

    it('should return empty array when no movements exist', () => {
      const transaction: Transaction = {
        id: '1',
        date: '2024-01-01',
        movements: [],
      };

      const result = TransactionUtils.getInputAccountIds(transaction);
      expect(result).toEqual([]);
    });

    it('should filter out movements without accountId', () => {
      const transaction: Transaction = {
        id: '1',
        date: '2024-01-01',
        movements: [
          { accountId: 'acc1', amount: -50, currencyId: 'usd' },
          { amount: -30, currencyId: 'usd' },
          { accountId: 'acc2', amount: 80, currencyId: 'usd' },
        ],
      };

      const result = TransactionUtils.getInputAccountIds(transaction);
      expect(result).toEqual(['acc1']);
    });
  });

  describe('getOutputAccountId', () => {
    it('should return the account ID with positive amount', () => {
      const transaction: Transaction = {
        id: '1',
        date: '2024-01-01',
        movements: [
          { accountId: 'acc1', amount: -100, currencyId: 'usd' },
          { accountId: 'acc2', amount: 100, currencyId: 'usd' },
        ],
      };

      const result = TransactionUtils.getOutputAccountId(transaction);
      expect(result).toBe('acc2');
    });

    it('should return undefined when no movements exist', () => {
      const transaction: Transaction = {
        id: '1',
        date: '2024-01-01',
        movements: [],
      };

      const result = TransactionUtils.getOutputAccountId(transaction);
      expect(result).toBeUndefined();
    });

    it('should return undefined when no positive amounts exist', () => {
      const transaction: Transaction = {
        id: '1',
        date: '2024-01-01',
        movements: [
          { accountId: 'acc1', amount: -100, currencyId: 'usd' },
          { accountId: 'acc2', amount: -50, currencyId: 'usd' },
        ],
      };

      const result = TransactionUtils.getOutputAccountId(transaction);
      expect(result).toBeUndefined();
    });

    it('should return the first output account when multiple exist', () => {
      const transaction: Transaction = {
        id: '1',
        date: '2024-01-01',
        movements: [
          { accountId: 'acc1', amount: -100, currencyId: 'usd' },
          { accountId: 'acc2', amount: 50, currencyId: 'usd' },
          { accountId: 'acc3', amount: 50, currencyId: 'usd' },
        ],
      };

      const result = TransactionUtils.getOutputAccountId(transaction);
      expect(result).toBe('acc2');
    });
  });

  describe('getOutputAccountIds', () => {
    it('should return all account IDs with positive amounts', () => {
      const transaction: Transaction = {
        id: '1',
        date: '2024-01-01',
        movements: [
          { accountId: 'acc1', amount: -100, currencyId: 'usd' },
          { accountId: 'acc2', amount: 50, currencyId: 'usd' },
          { accountId: 'acc3', amount: 50, currencyId: 'usd' },
        ],
      };

      const result = TransactionUtils.getOutputAccountIds(transaction);
      expect(result).toEqual(['acc2', 'acc3']);
    });

    it('should return empty array when no movements exist', () => {
      const transaction: Transaction = {
        id: '1',
        date: '2024-01-01',
        movements: [],
      };

      const result = TransactionUtils.getOutputAccountIds(transaction);
      expect(result).toEqual([]);
    });

    it('should filter out movements without accountId', () => {
      const transaction: Transaction = {
        id: '1',
        date: '2024-01-01',
        movements: [
          { accountId: 'acc1', amount: -80, currencyId: 'usd' },
          { accountId: 'acc2', amount: 50, currencyId: 'usd' },
          { amount: 30, currencyId: 'usd' },
        ],
      };

      const result = TransactionUtils.getOutputAccountIds(transaction);
      expect(result).toEqual(['acc2']);
    });
  });
});

