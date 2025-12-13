import { Transaction } from '../../../core/api/models/transaction';

/**
 * Represents the effective amount for a specific currency in a transaction
 */
export interface EffectiveAmount {
  currencyId: string;
  amount: number;
}

/**
 * Utility class for Transaction operations
 */
export class TransactionUtils {
  /**
   * Get the input account ID (source of money - negative amount in movements)
   * @param transaction The transaction to analyze
   * @returns The account ID that is the source of money, or undefined if not found
   */
  static getInputAccountId(transaction: Transaction): string | undefined {
    if (!transaction.movements || transaction.movements.length === 0) {
      return undefined;
    }

    // Find the first movement with negative amount (money going out)
    const inputMovement = transaction.movements.find(movement => movement.amount < 0);
    return inputMovement?.accountId;
  }

  /**
   * Get all input account IDs (sources of money - negative amounts in movements)
   * @param transaction The transaction to analyze
   * @returns Array of account IDs that are sources of money
   */
  static getInputAccountIds(transaction: Transaction): string[] {
    if (!transaction.movements || transaction.movements.length === 0) {
      return [];
    }

    // Find all movements with negative amounts (money going out)
    return transaction.movements
      .filter(movement => movement.amount < 0 && movement.accountId)
      .map(movement => movement.accountId!);
  }

  /**
   * Get the output account ID (destination of money - positive amount in movements)
   * @param transaction The transaction to analyze
   * @returns The account ID that is the destination of money, or undefined if not found
   */
  static getOutputAccountId(transaction: Transaction): string | undefined {
    if (!transaction.movements || transaction.movements.length === 0) {
      return undefined;
    }

    // Find the first movement with positive amount (money coming in)
    const outputMovement = transaction.movements.find(movement => movement.amount > 0);
    return outputMovement?.accountId;
  }

  /**
   * Get all output account IDs (destinations of money - positive amounts in movements)
   * @param transaction The transaction to analyze
   * @returns Array of account IDs that are destinations of money
   */
  static getOutputAccountIds(transaction: Transaction): string[] {
    if (!transaction.movements || transaction.movements.length === 0) {
      return [];
    }

    // Find all movements with positive amounts (money coming in)
    return transaction.movements
      .filter(movement => movement.amount > 0 && movement.accountId)
      .map(movement => movement.accountId!);
  }

  /**
   * Calculate the effective amount per currency for a transaction.
   * The effective amount represents how much money of each currency was actually used.
   *
   * For each currency:
   * - Group movements by currency
   * - Calculate the sum of positive amounts (credits)
   * - Calculate the sum of negative amounts (debits)
   * - The effective amount is the maximum of the absolute values of these sums
   *
   * Example: If a transaction has movements of +100 USD, -50 USD, and -30 USD,
   * the effective amount would be max(abs(100), abs(-50 + -30)) = max(100, 80) = 100 USD
   *
   * @param transaction The transaction to analyze
   * @returns Array of effective amounts per currency
   */
  static getEffectiveAmounts(transaction: Transaction): EffectiveAmount[] {
    if (!transaction.movements || transaction.movements.length === 0) {
      return [];
    }

    // Group movements by currency
    const currencyGroups = new Map<string, { positive: number; negative: number }>();

    transaction.movements.forEach(movement => {
      const currencyId = movement.currencyId;
      if (!currencyGroups.has(currencyId)) {
        currencyGroups.set(currencyId, { positive: 0, negative: 0 });
      }

      const group = currencyGroups.get(currencyId)!;
      if (movement.amount > 0) {
        group.positive += movement.amount;
      } else {
        group.negative += movement.amount;
      }
    });

    // Calculate effective amount for each currency
    const effectiveAmounts: EffectiveAmount[] = [];
    currencyGroups.forEach((group, currencyId) => {
      const effectiveAmount = Math.max(
        Math.abs(group.positive),
        Math.abs(group.negative)
      );
      effectiveAmounts.push({
        currencyId,
        amount: effectiveAmount
      });
    });

    return effectiveAmounts;
  }
}

