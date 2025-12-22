import { Routes } from '@angular/router';
import { authGuard } from './core/auth/guards/auth.guard';
import { noAuthGuard } from './core/auth/guards/no-auth.guard';
import { LayoutComponent } from './layout/layout.component';

export const routes: Routes = [
    {
        path: 'auth',
        canActivate: [noAuthGuard],
        children: [
            {
                path: 'login',
                loadComponent: () =>
                    import('./features/auth/login/login.component').then((m) => m.LoginComponent),
            },
            {
                path: '',
                redirectTo: 'login',
                pathMatch: 'full',
            },
        ],
    },
    {
        path: '',
        component: LayoutComponent,
        canActivate: [authGuard],
        children: [
            {
                path: 'dashboard',
                loadComponent: () =>
                    import('./features/dashboard/dashboard.component').then(
                        (m) => m.DashboardComponent,
                    ),
            },
            {
                path: 'transactions',
                loadComponent: () =>
                    import('./features/transactions/transactions.component').then(
                        (m) => m.TransactionsComponent,
                    ),
            },
            {
                path: 'suspicious',
                loadComponent: () =>
                    import('./features/transactions/suspicious-transactions/suspicious-transactions.component').then(
                        (m) => m.SuspiciousTransactionsComponent,
                    ),
            },
            {
                path: 'accounts',
                loadComponent: () =>
                    import('./features/accounts/accounts.component').then(
                        (m) => m.AccountsComponent,
                    ),
            },
            {
                path: 'currencies',
                loadComponent: () =>
                    import('./features/currencies/currencies.component').then(
                        (m) => m.CurrenciesComponent,
                    ),
            },
            {
                path: 'bank-importers',
                loadComponent: () =>
                    import('./features/bank-importers/bank-importers.component').then(
                        (m) => m.BankImportersComponent,
                    ),
            },
            {
                path: 'unprocessed',
                loadComponent: () =>
                    import('./features/unprocessed-transactions/unprocessed-transactions.component').then(
                        (m) => m.UnprocessedTransactionsComponent,
                    ),
            },
            {
                path: 'matchers',
                loadComponent: () =>
                    import('./features/matchers/matchers.component').then(
                        (m) => m.MatchersComponent,
                    ),
            },
            {
                path: 'budget',
                loadComponent: () =>
                    import('./features/budget-items/budget-items.component').then(
                        (m) => m.BudgetItemsComponent,
                    ),
            },
            {
                path: 'reports',
                loadComponent: () =>
                    import('./features/reports/reports.component').then((m) => m.ReportsComponent),
                children: [
                    {
                        path: 'expense',
                        loadComponent: () =>
                            import('./features/reports/expense-report/expense-report.component').then(
                                (m) => m.ExpenseReportComponent,
                            ),
                    },
                    {
                        path: 'balance',
                        loadComponent: () =>
                            import('./features/reports/balance-report/balance-report.component').then(
                                (m) => m.BalanceReportComponent,
                            ),
                    },
                ],
            },
            {
                path: 'notifications',
                loadComponent: () =>
                    import('./features/notifications/notifications.component').then(
                        (m) => m.NotificationsComponent,
                    ),
            },
            {
                path: 'settings',
                loadComponent: () =>
                    import('./features/settings/settings.component').then(
                        (m) => m.SettingsComponent,
                    ),
            },
            {
                path: '',
                redirectTo: 'dashboard',
                pathMatch: 'full',
            },
        ],
    },
    {
        path: '**',
        redirectTo: 'dashboard',
    },
];
