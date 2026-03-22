import { Component, inject, OnInit } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatTableModule } from '@angular/material/table';
import { MatDialog, MatDialogModule } from '@angular/material/dialog';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TemplateService } from './services/template.service';
import { AccountService } from '../accounts/services/account.service';
import { TemplateEditDialogComponent } from './template-edit-dialog/template-edit-dialog.component';
import { TransactionTemplate } from '../../core/api/models/transaction-template';

@Component({
    selector: 'app-templates',
    standalone: true,
    imports: [
        MatButtonModule,
        MatIconModule,
        MatTableModule,
        MatDialogModule,
        MatProgressSpinnerModule,
        MatTooltipModule,
    ],
    templateUrl: './templates.component.html',
})
export class TemplatesComponent implements OnInit {
    private readonly templateService = inject(TemplateService);
    private readonly accountService = inject(AccountService);
    private readonly dialog = inject(MatDialog);

    protected readonly templates = this.templateService.templates;
    protected readonly loading = this.templateService.loading;
    protected readonly displayedColumns = ['name', 'description', 'accounts', 'actions'];

    ngOnInit(): void {
        this.accountService.loadAccounts().subscribe();
        this.templateService.loadTemplates().subscribe();
    }

    protected getAccountNames(template: TransactionTemplate): string {
        const accounts = this.accountService.accounts();
        const accountMap = new Map(accounts.map((a) => [a.id, a.name]));
        const ids = [
            ...new Set((template.movements ?? []).map((m) => m.accountId).filter(Boolean)),
        ];
        return ids.map((id) => accountMap.get(id!) ?? id!).join(', ');
    }

    protected openCreateDialog(): void {
        this.dialog.open(TemplateEditDialogComponent, {
            width: '640px',
            data: {},
            disableClose: true,
        });
    }

    protected openEditDialog(template: TransactionTemplate): void {
        this.dialog.open(TemplateEditDialogComponent, {
            width: '640px',
            data: { template },
            disableClose: true,
        });
    }

    protected delete(template: TransactionTemplate): void {
        if (confirm(`Delete template "${template.name}"?`)) {
            this.templateService.delete(template.id).subscribe();
        }
    }
}
