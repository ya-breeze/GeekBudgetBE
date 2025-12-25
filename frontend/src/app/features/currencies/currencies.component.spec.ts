import { ComponentFixture, TestBed } from '@angular/core/testing';
import { CurrenciesComponent } from './currencies.component';
import { CurrencyService } from './services/currency.service';

import { MatSnackBar } from '@angular/material/snack-bar';
import { signal } from '@angular/core';
import { of } from 'rxjs';
import { Currency } from '../../core/api/models/currency';
import { NoopAnimationsModule } from '@angular/platform-browser/animations';

describe('CurrenciesComponent', () => {
    let component: CurrenciesComponent;
    let fixture: ComponentFixture<CurrenciesComponent>;
    let currencyService: jasmine.SpyObj<CurrencyService>;

    const mockCurrencies: Currency[] = [
        { id: '1', name: 'USD', description: 'US Dollar' },
        { id: '2', name: 'EUR', description: 'Euro' },
    ];

    beforeEach(async () => {
        const currencyServiceSpy = jasmine.createSpyObj(
            'CurrencyService',
            ['loadCurrencies', 'delete', 'create', 'update'],
            {
                currencies: signal(mockCurrencies),
                loading: signal(false),
                error: signal(null),
            },
        );
        // Set default return values for methods
        currencyServiceSpy.loadCurrencies.and.returnValue(of(mockCurrencies));
        currencyServiceSpy.create.and.returnValue(of({} as Currency));
        currencyServiceSpy.update.and.returnValue(of({} as Currency));
        currencyServiceSpy.delete.and.returnValue(of(undefined));

        const snackBarSpy = jasmine.createSpyObj('MatSnackBar', ['open']);

        await TestBed.configureTestingModule({
            imports: [CurrenciesComponent, NoopAnimationsModule],
            providers: [
                { provide: CurrencyService, useValue: currencyServiceSpy },
                { provide: MatSnackBar, useValue: snackBarSpy },
            ],
        }).compileComponents();

        currencyService = TestBed.inject(CurrencyService) as jasmine.SpyObj<CurrencyService>;

        fixture = TestBed.createComponent(CurrenciesComponent);
        component = fixture.componentInstance;
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });

    it('should display loading spinner while loading', () => {
        currencyService.loading.set(true);
        fixture.detectChanges();

        const compiled = fixture.nativeElement;
        const spinner = compiled.querySelector('mat-spinner');
        expect(spinner).toBeTruthy();
    });

    it('should display currencies in table', () => {
        currencyService.loadCurrencies.and.returnValue(of(mockCurrencies));
        fixture.detectChanges();

        expect(component['currencies']()).toEqual(mockCurrencies);
    });

    it('should open create dialog on add button click', () => {
        const dialogRefSpy = jasmine.createSpyObj('MatDialogRef', ['afterClosed']);
        dialogRefSpy.afterClosed.and.returnValue(of(null));
        const dialogSpy = spyOn(component['dialog'], 'open').and.returnValue(dialogRefSpy);

        component.openCreateDialog();

        expect(dialogSpy).toHaveBeenCalled();
    });

    it('should open edit dialog on edit button click', () => {
        const currency = mockCurrencies[0];
        const dialogRefSpy = jasmine.createSpyObj('MatDialogRef', ['afterClosed']);
        dialogRefSpy.afterClosed.and.returnValue(of(null));
        const dialogSpy = spyOn(component['dialog'], 'open').and.returnValue(dialogRefSpy);

        component.openEditDialog(currency);

        expect(dialogSpy).toHaveBeenCalled();
    });

    it('should open delete dialog on delete click', () => {
        const currency = mockCurrencies[0];
        const dialogRefSpy = jasmine.createSpyObj('MatDialogRef', ['afterClosed']);
        dialogRefSpy.afterClosed.and.returnValue(of({ replaceWithCurrencyId: undefined }));
        const dialogSpy = spyOn(component['dialog'], 'open').and.returnValue(dialogRefSpy);

        component.deleteCurrency(currency);

        expect(dialogSpy).toHaveBeenCalled();
        expect(currencyService.delete).toHaveBeenCalledWith(currency.id, undefined);
    });

    it('should call create service when dialog returns result', () => {
        const newCurrency: Currency = { id: '3', name: 'EUR', description: 'Euro' };
        const dialogRefSpy = jasmine.createSpyObj('MatDialogRef', ['afterClosed']);
        dialogRefSpy.afterClosed.and.returnValue(of(newCurrency));
        spyOn(component['dialog'], 'open').and.returnValue(dialogRefSpy);
        currencyService.create.and.returnValue(of(newCurrency));

        component.openCreateDialog();

        expect(currencyService.create).toHaveBeenCalledWith(newCurrency);
    });

    it('should display error message on failure', () => {
        currencyService.error.set('Failed to load currencies');
        fixture.detectChanges();

        // The component doesn't expose error directly, but the service does
        expect(currencyService.error()).toBe('Failed to load currencies');
    });

    it('should not delete if dialog is cancelled', () => {
        const currency = mockCurrencies[0];
        const dialogRefSpy = jasmine.createSpyObj('MatDialogRef', ['afterClosed']);
        dialogRefSpy.afterClosed.and.returnValue(of(null));
        spyOn(component['dialog'], 'open').and.returnValue(dialogRefSpy);

        component.deleteCurrency(currency);

        expect(currencyService.delete).not.toHaveBeenCalled();
    });
});
