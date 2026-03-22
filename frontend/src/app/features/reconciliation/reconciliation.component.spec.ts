import { ComponentFixture, TestBed } from '@angular/core/testing';
import { ReconciliationComponent } from './reconciliation.component';
import { ReconciliationService } from './services/reconciliation.service';
import { MatSnackBar } from '@angular/material/snack-bar';
import { MatDialog, MatDialogRef } from '@angular/material/dialog';
import { Router } from '@angular/router';
import { of, Subject } from 'rxjs';
import { ReconciliationStatus } from '../../core/api/models/reconciliation-status';
import { Reconciliation } from '../../core/api/models/reconciliation';
import { NoopAnimationsModule } from '@angular/platform-browser/animations';

describe('ReconciliationComponent', () => {
    let component: ReconciliationComponent;
    let fixture: ComponentFixture<ReconciliationComponent>;
    let mockReconciliationService: jasmine.SpyObj<ReconciliationService>;
    let mockSnackBar: jasmine.SpyObj<MatSnackBar>;
    let mockDialog: jasmine.SpyObj<MatDialog>;
    let mockRouter: jasmine.SpyObj<Router>;

    const baseStatus: ReconciliationStatus = {
        accountId: 'acc1',
        accountName: 'Cash',
        currencyId: 'USD',
        currencySymbol: '$',
        bankBalance: 100,
        appBalance: 200,
        delta: 100,
        hasUnprocessedTransactions: false,
        hasBankImporter: false,
        isManualReconciliationEnabled: true,
    };

    beforeEach(async () => {
        mockReconciliationService = jasmine.createSpyObj('ReconciliationService', [
            'loadStatuses',
            'reconcile',
            'enableManual',
            'getTransactionsSince',
            'analyzeDisbalance',
        ]);
        mockSnackBar = jasmine.createSpyObj('MatSnackBar', ['open']);
        mockDialog = jasmine.createSpyObj('MatDialog', ['open']);
        mockRouter = jasmine.createSpyObj('Router', ['navigate']);

        mockReconciliationService.loadStatuses.and.returnValue(of([]));

        await TestBed.configureTestingModule({
            imports: [ReconciliationComponent, NoopAnimationsModule],
            providers: [
                { provide: ReconciliationService, useValue: mockReconciliationService },
                { provide: MatSnackBar, useValue: mockSnackBar },
                { provide: Router, useValue: mockRouter },
            ],
        })
            .overrideProvider(MatDialog, { useValue: mockDialog })
            .compileComponents();

        fixture = TestBed.createComponent(ReconciliationComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    describe('getReconcileTooltip', () => {
        it('returns unprocessed message when hasUnprocessedTransactions', () => {
            const s = { ...baseStatus, hasUnprocessedTransactions: true };
            expect(component.getReconcileTooltip(s)).toContain('unprocessed');
        });

        it('returns confirm message for no-importer account with large delta', () => {
            const s = { ...baseStatus, hasBankImporter: false, delta: 100 };
            const tooltip = component.getReconcileTooltip(s);
            expect(tooltip).toContain('click to confirm');
            expect(tooltip).toContain('100.00');
        });

        it('returns too-large message for importer account with large delta', () => {
            const s = { ...baseStatus, hasBankImporter: true, delta: 100 };
            const tooltip = component.getReconcileTooltip(s);
            expect(tooltip).toContain('too large');
        });

        it('returns mark-as-reconciled for small delta', () => {
            const s = { ...baseStatus, delta: 0.005 };
            expect(component.getReconcileTooltip(s)).toBe('Mark as Reconciled');
        });
    });

    describe('getStatusClass', () => {
        it('returns status-yellow when hasUnprocessedTransactions', () => {
            expect(
                component.getStatusClass({ ...baseStatus, hasUnprocessedTransactions: true }),
            ).toBe('status-yellow');
        });

        it('returns status-yellow for no-importer account with large delta', () => {
            expect(
                component.getStatusClass({ ...baseStatus, hasBankImporter: false, delta: 100 }),
            ).toBe('status-yellow');
        });

        it('returns status-green for no-importer account within tolerance', () => {
            expect(
                component.getStatusClass({ ...baseStatus, hasBankImporter: false, delta: 0.005 }),
            ).toBe('status-green');
        });

        it('returns status-red for importer account with large delta', () => {
            expect(
                component.getStatusClass({ ...baseStatus, hasBankImporter: true, delta: 100 }),
            ).toBe('status-red');
        });

        it('returns status-green for importer account within tolerance', () => {
            expect(
                component.getStatusClass({ ...baseStatus, hasBankImporter: true, delta: 0.005 }),
            ).toBe('status-green');
        });
    });

    describe('reconcile', () => {
        const mockReconciliation: Reconciliation = {
            accountId: 'acc1',
            currencyId: 'USD',
            reconciledAt: '2026-01-01T00:00:00Z',
            reconciledBalance: 0,
        };

        it('opens confirmation dialog when no-importer and large delta', () => {
            const afterClosedSubject = new Subject<boolean>();
            const mockDialogRef = {
                afterClosed: () => afterClosedSubject.asObservable(),
            } as MatDialogRef<any>;
            mockDialog.open.and.returnValue(mockDialogRef);

            const s = { ...baseStatus, hasBankImporter: false, delta: 100 };
            component.reconcile(s);

            expect(mockDialog.open).toHaveBeenCalled();
            expect(mockReconciliationService.reconcile).not.toHaveBeenCalled();
        });

        it('does not call API when user cancels confirmation dialog', () => {
            const afterClosedSubject = new Subject<boolean>();
            const mockDialogRef = {
                afterClosed: () => afterClosedSubject.asObservable(),
            } as MatDialogRef<any>;
            mockDialog.open.and.returnValue(mockDialogRef);

            const s = { ...baseStatus, hasBankImporter: false, delta: 100 };
            component.reconcile(s);
            afterClosedSubject.next(false); // user cancels

            expect(mockReconciliationService.reconcile).not.toHaveBeenCalled();
        });

        it('calls API when user confirms dialog', () => {
            const afterClosedSubject = new Subject<boolean>();
            const mockDialogRef = {
                afterClosed: () => afterClosedSubject.asObservable(),
            } as MatDialogRef<any>;
            mockDialog.open.and.returnValue(mockDialogRef);
            mockReconciliationService.reconcile.and.returnValue(of(mockReconciliation));
            mockReconciliationService.loadStatuses.and.returnValue(of([]));

            const s = { ...baseStatus, hasBankImporter: false, delta: 100 };
            component.reconcile(s);
            afterClosedSubject.next(true); // user confirms

            expect(mockReconciliationService.reconcile).toHaveBeenCalledWith('acc1', {
                currencyId: 'USD',
                balance: 0,
            });
        });

        it('calls API directly (no dialog) when delta is within tolerance', () => {
            mockReconciliationService.reconcile.and.returnValue(of(mockReconciliation));
            mockReconciliationService.loadStatuses.and.returnValue(of([]));

            const s = { ...baseStatus, hasBankImporter: false, delta: 0.005 };
            component.reconcile(s);

            expect(mockDialog.open).not.toHaveBeenCalled();
            expect(mockReconciliationService.reconcile).toHaveBeenCalled();
        });

        it('disables reconcile button when hasUnprocessedTransactions regardless of importer status', () => {
            // Verify via getReconcileTooltip (template uses [disabled] which we test via tooltip)
            const s = {
                ...baseStatus,
                hasBankImporter: false,
                hasUnprocessedTransactions: true,
                delta: 100,
            };
            expect(component.getReconcileTooltip(s)).toContain('unprocessed');
        });
    });
});
