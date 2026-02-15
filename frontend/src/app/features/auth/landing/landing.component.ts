import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';

@Component({
    selector: 'app-landing',
    standalone: true,
    imports: [CommonModule, RouterModule, MatCardModule, MatButtonModule, MatIconModule],
    templateUrl: './landing.component.html',
    styleUrls: ['./landing.component.scss'],
})
export class LandingComponent {
    features = [
        {
            icon: 'cloud_download',
            title: 'Bank Import',
            description: 'Import transactions from FIO, Komerční banka, and Revolut',
            details:
                'Automatically sync your bank transactions with support for multiple currencies and accounts. No manual data entry required.',
        },
        {
            icon: 'bolt',
            title: 'Smart Matching',
            description: 'Automated transaction categorization and partner recognition',
            details:
                'Create custom matchers to automatically categorize transactions based on patterns, amounts, and partners. Save time with intelligent automation.',
        },
        {
            icon: 'track_changes',
            title: 'Budget Planning',
            description: 'Set targets and track spending by category',
            details:
                'Plan your monthly budgets with flexible target amounts per category. Monitor your progress and stay within your financial goals.',
        },
        {
            icon: 'check_circle',
            title: 'Reconciliation',
            description: 'Balance tracking and account verification',
            details:
                'Verify your account balances match your bank statements. Track discrepancies and ensure your financial data is always accurate.',
        },
        {
            icon: 'merge',
            title: 'Duplicate Detection',
            description: 'Automatic identification and merging of duplicates',
            details:
                'Detect and merge duplicate transactions automatically. Keep your financial records clean and accurate without manual cleanup.',
        },
        {
            icon: 'bar_chart',
            title: 'Analytics',
            description: 'Visualize spending patterns and trends',
            details:
                'Track your financial health with comprehensive dashboards, charts, and reports. Make informed decisions based on your spending patterns.',
        },
    ];

    steps = [
        {
            number: 1,
            title: 'Import Transactions',
            description:
                'Connect your bank accounts and import transactions automatically or upload bank statements',
        },
        {
            number: 2,
            title: 'Review & Match',
            description:
                'Process unprocessed transactions and set up matchers to automatically categorize future transactions',
        },
        {
            number: 3,
            title: 'Set Budgets',
            description:
                'Create monthly budgets for different categories and track your spending against targets',
        },
        {
            number: 4,
            title: 'Monitor & Optimize',
            description:
                'Track your progress with dashboards, reconcile accounts, and adjust your financial strategy',
        },
    ];

    concepts = [
        {
            icon: 'filter_alt',
            title: 'Matchers',
            description:
                'Automated rules that categorize transactions based on patterns like partner name, amount range, or transaction type. Create matchers once and let GeekBudget automatically process future transactions.',
        },
        {
            icon: 'pending_actions',
            title: 'Unprocessed Transactions',
            description:
                "Imported transactions that haven't been categorized yet. Review these transactions to assign categories and create matchers for similar future transactions.",
        },
        {
            icon: 'track_changes',
            title: 'Budget',
            description:
                'Monthly spending targets for different categories. Set realistic goals and track your actual spending against planned amounts to stay on track.',
        },
        {
            icon: 'event_available',
            title: 'Reconciliation',
            description:
                'The process of verifying your account balance in GeekBudget matches your actual bank balance. Helps identify missing or duplicate transactions.',
        },
        {
            icon: 'currency_exchange',
            title: 'Multi-Currency Support',
            description:
                'Track transactions in different currencies with automatic conversion. Perfect for managing multiple accounts across different countries.',
        },
        {
            icon: 'category',
            title: 'Categories',
            description:
                'Organize your spending into meaningful groups like "Groceries", "Transportation", or "Entertainment". Use categories for budgeting and financial analysis.',
        },
    ];

    scrollToFeatures(): void {
        const element = document.getElementById('features');
        element?.scrollIntoView({ behavior: 'smooth' });
    }
}
