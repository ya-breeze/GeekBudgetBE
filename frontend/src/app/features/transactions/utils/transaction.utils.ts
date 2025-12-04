import { Transaction } from '../../../core/api/models/transaction';

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
}

