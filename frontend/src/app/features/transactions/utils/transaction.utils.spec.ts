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

  describe('getEffectiveAmounts', () => {
    it('should return empty array when no movements exist', () => {
      const transaction: Transaction = {
        id: '1',
        date: '2024-01-01',
        movements: [],
      };

      const result = TransactionUtils.getEffectiveAmounts(transaction);
      expect(result).toEqual([]);
    });

    it('should calculate effective amount for simple single currency transaction', () => {
      const transaction: Transaction = {
        id: '1',
        date: '2024-01-01',
        movements: [
          { accountId: 'acc1', amount: -100, currencyId: 'usd' },
          { accountId: 'acc2', amount: 100, currencyId: 'usd' },
        ],
      };

      const result = TransactionUtils.getEffectiveAmounts(transaction);
      expect(result.length).toBe(1);
      expect(result[0].currencyId).toBe('usd');
      expect(result[0].amount).toBe(100);
    });

    it('should calculate effective amount as max of positive and negative sums', () => {
      // Example from requirements: +100 USD, -50 USD, -30 USD
      // Effective amount = max(abs(100), abs(-50 + -30)) = max(100, 80) = 100
      const transaction: Transaction = {
        id: '1',
        date: '2024-01-01',
        movements: [
          { accountId: 'acc1', amount: 100, currencyId: 'usd' },
          { accountId: 'acc2', amount: -50, currencyId: 'usd' },
          { accountId: 'acc3', amount: -30, currencyId: 'usd' },
        ],
      };

      const result = TransactionUtils.getEffectiveAmounts(transaction);
      expect(result.length).toBe(1);
      expect(result[0].currencyId).toBe('usd');
      expect(result[0].amount).toBe(100);
    });

    it('should calculate effective amount when negative sum is larger', () => {
      // +50 USD, -80 USD, -30 USD
      // Effective amount = max(abs(50), abs(-80 + -30)) = max(50, 110) = 110
      const transaction: Transaction = {
        id: '1',
        date: '2024-01-01',
        movements: [
          { accountId: 'acc1', amount: 50, currencyId: 'usd' },
          { accountId: 'acc2', amount: -80, currencyId: 'usd' },
          { accountId: 'acc3', amount: -30, currencyId: 'usd' },
        ],
      };

      const result = TransactionUtils.getEffectiveAmounts(transaction);
      expect(result.length).toBe(1);
      expect(result[0].currencyId).toBe('usd');
      expect(result[0].amount).toBe(110);
    });

    it('should handle multiple currencies independently', () => {
      const transaction: Transaction = {
        id: '1',
        date: '2024-01-01',
        movements: [
          { accountId: 'acc1', amount: -100, currencyId: 'usd' },
          { accountId: 'acc2', amount: 100, currencyId: 'usd' },
          { accountId: 'acc3', amount: -50, currencyId: 'eur' },
          { accountId: 'acc4', amount: 50, currencyId: 'eur' },
        ],
      };

      const result = TransactionUtils.getEffectiveAmounts(transaction);
      expect(result.length).toBe(2);

      const usdAmount = result.find(ea => ea.currencyId === 'usd');
      expect(usdAmount).toBeDefined();
      expect(usdAmount!.amount).toBe(100);

      const eurAmount = result.find(ea => ea.currencyId === 'eur');
      expect(eurAmount).toBeDefined();
      expect(eurAmount!.amount).toBe(50);
    });

    it('should handle complex multi-currency transaction', () => {
      // USD: +200, -100, -50 => max(200, 150) = 200
      // EUR: +80, -100 => max(80, 100) = 100
      // GBP: -75, +75 => max(75, 75) = 75
      const transaction: Transaction = {
        id: '1',
        date: '2024-01-01',
        movements: [
          { accountId: 'acc1', amount: 200, currencyId: 'usd' },
          { accountId: 'acc2', amount: -100, currencyId: 'usd' },
          { accountId: 'acc3', amount: -50, currencyId: 'usd' },
          { accountId: 'acc4', amount: 80, currencyId: 'eur' },
          { accountId: 'acc5', amount: -100, currencyId: 'eur' },
          { accountId: 'acc6', amount: -75, currencyId: 'gbp' },
          { accountId: 'acc7', amount: 75, currencyId: 'gbp' },
        ],
      };

      const result = TransactionUtils.getEffectiveAmounts(transaction);
      expect(result.length).toBe(3);

      const usdAmount = result.find(ea => ea.currencyId === 'usd');
      expect(usdAmount!.amount).toBe(200);

      const eurAmount = result.find(ea => ea.currencyId === 'eur');
      expect(eurAmount!.amount).toBe(100);

      const gbpAmount = result.find(ea => ea.currencyId === 'gbp');
      expect(gbpAmount!.amount).toBe(75);
    });

    it('should handle transaction with only positive amounts', () => {
      const transaction: Transaction = {
        id: '1',
        date: '2024-01-01',
        movements: [
          { accountId: 'acc1', amount: 100, currencyId: 'usd' },
          { accountId: 'acc2', amount: 50, currencyId: 'usd' },
        ],
      };

      const result = TransactionUtils.getEffectiveAmounts(transaction);
      expect(result.length).toBe(1);
      expect(result[0].currencyId).toBe('usd');
      expect(result[0].amount).toBe(150); // max(150, 0) = 150
    });

    it('should handle transaction with only negative amounts', () => {
      const transaction: Transaction = {
        id: '1',
        date: '2024-01-01',
        movements: [
          { accountId: 'acc1', amount: -100, currencyId: 'usd' },
          { accountId: 'acc2', amount: -50, currencyId: 'usd' },
        ],
      };

      const result = TransactionUtils.getEffectiveAmounts(transaction);
      expect(result.length).toBe(1);
      expect(result[0].currencyId).toBe('usd');
      expect(result[0].amount).toBe(150); // max(0, 150) = 150
    });

    it('should handle movements without accountId', () => {
      const transaction: Transaction = {
        id: '1',
        date: '2024-01-01',
        movements: [
          { amount: -100, currencyId: 'usd' },
          { amount: 100, currencyId: 'usd' },
        ],
      };

      const result = TransactionUtils.getEffectiveAmounts(transaction);
      expect(result.length).toBe(1);
      expect(result[0].currencyId).toBe('usd');
      expect(result[0].amount).toBe(100);
    });

    it('should handle decimal amounts correctly', () => {
      const transaction: Transaction = {
        id: '1',
        date: '2024-01-01',
        movements: [
          { accountId: 'acc1', amount: 123.45, currencyId: 'usd' },
          { accountId: 'acc2', amount: -67.89, currencyId: 'usd' },
          { accountId: 'acc3', amount: -55.56, currencyId: 'usd' },
        ],
      };

      const result = TransactionUtils.getEffectiveAmounts(transaction);
      expect(result.length).toBe(1);
      expect(result[0].currencyId).toBe('usd');
      expect(result[0].amount).toBeCloseTo(123.45, 2);
    });
  });
});

