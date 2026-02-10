import { Injectable, signal, computed } from '@angular/core';
import { Transaction } from '../../../core/api/models/transaction';

@Injectable({
    providedIn: 'root',
})
export class TransactionSelectionService {
    private readonly selectedTransactionsSignal = signal<Transaction[]>([]);

    readonly selectedTransactions = computed(() => this.selectedTransactionsSignal());
    readonly count = computed(() => this.selectedTransactionsSignal().length);

    toggleSelection(transaction: Transaction): void {
        const current = this.selectedTransactionsSignal();
        const index = current.findIndex((t) => t.id === transaction.id);

        if (index !== -1) {
            // Remove if already selected
            this.selectedTransactionsSignal.set(current.filter((t) => t.id !== transaction.id));
        } else {
            // Add if not selected, limit to 2
            if (current.length < 2) {
                this.selectedTransactionsSignal.set([...current, transaction]);
            }
        }
    }

    isSelected(transactionId: string | undefined): boolean {
        if (!transactionId) return false;
        return this.selectedTransactionsSignal().some((t) => t.id === transactionId);
    }

    clearSelection(): void {
        this.selectedTransactionsSignal.set([]);
    }

    remove(transactionId: string): void {
        const current = this.selectedTransactionsSignal();
        this.selectedTransactionsSignal.set(current.filter((t) => t.id !== transactionId));
    }
}
