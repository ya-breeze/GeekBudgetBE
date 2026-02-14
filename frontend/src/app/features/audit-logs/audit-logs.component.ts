import { Component, inject, OnInit, ViewChild, signal, AfterViewInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { MatTableModule, MatTableDataSource } from '@angular/material/table';
import { MatSort, MatSortModule } from '@angular/material/sort';
import { MatPaginator, MatPaginatorModule } from '@angular/material/paginator';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { MatNativeDateModule } from '@angular/material/core';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatDialog, MatDialogModule } from '@angular/material/dialog';
import { AuditLogsService } from './services/audit-logs.service';
import { AuditLog } from '../../core/api/models/audit-log';
import { AppDatePipe } from '../../shared/pipes/app-date.pipe';
import { AuditLogDetailDialogComponent } from './audit-log-detail-dialog/audit-log-detail-dialog.component';

@Component({
    selector: 'app-audit-logs',
    standalone: true,
    imports: [
        CommonModule,
        FormsModule,
        MatTableModule,
        MatSortModule,
        MatPaginatorModule,
        MatButtonModule,
        MatIconModule,
        MatProgressSpinnerModule,
        MatFormFieldModule,
        MatInputModule,
        MatDatepickerModule,
        MatNativeDateModule,
        MatTooltipModule,
        MatDialogModule,
        AppDatePipe,
    ],
    templateUrl: './audit-logs.component.html',
    styleUrl: './audit-logs.component.scss',
})
export class AuditLogsComponent implements OnInit, AfterViewInit {
    private readonly auditLogsService = inject(AuditLogsService);
    private readonly dialog = inject(MatDialog);

    protected readonly auditLogs = this.auditLogsService.auditLogs;
    protected readonly loading = this.auditLogsService.loading;

    displayedColumns: string[] = [
        'createdAt',
        'action',
        'entityType',
        'entityId',
        'changeSource',
        'snapshot',
    ];
    dataSource = new MatTableDataSource<AuditLog>([]);

    // Filters
    entityType = signal('');
    entityId = signal('');
    dateFrom = signal<Date | null>(null);
    dateTo = signal<Date | null>(null);

    @ViewChild(MatSort) sort!: MatSort;
    @ViewChild(MatPaginator) paginator!: MatPaginator;

    ngOnInit(): void {
        this.loadData();
    }

    ngAfterViewInit() {
        this.dataSource.sort = this.sort;
        this.dataSource.paginator = this.paginator;

        // Update dataSource when signal changes
        // Effect isn't directly available in standard lifecycle without constructor or injection context,
        // but since auditLogs is a signal, we can just use computed or subscribe in effect if we migrated to effect()
        // Here we can simply subscribe or check changes.
        // For simplicity with signals in template, we might bind directly, but MatTableDataSource needs array.
        // So we'll use an effect or just update manually when load finishes.
    }

    // Using effect logic manually since we are in a class method
    constructor() {
        // We can use effect here if needed, but loadData handles the call.
        // Syncing signal to dataSource:
        // This is a bit manual. Alternatively use a computed/effect.
    }

    loadData(): void {
        const params: any = {};
        if (this.entityType()) params.entityType = this.entityType();
        if (this.entityId()) params.entityId = this.entityId();
        if (this.dateFrom()) params.dateFrom = this.dateFrom()!.toISOString();
        if (this.dateTo()) params.dateTo = this.dateTo()!.toISOString();

        this.auditLogsService.loadAuditLogs(params).subscribe((logs) => {
            this.dataSource.data = logs;
            if (this.paginator) {
                this.dataSource.paginator = this.paginator;
            }
            if (this.sort) {
                this.dataSource.sort = this.sort;
            }
        });
    }

    clearFilters(): void {
        this.entityType.set('');
        this.entityId.set('');
        this.dateFrom.set(null);
        this.dateTo.set(null);
        this.loadData();
    }

    openDetailDialog(log: AuditLog): void {
        this.dialog.open(AuditLogDetailDialogComponent, {
            data: log,
            width: '95vw',
            maxWidth: '95vw',
            height: '90vh',
        });
    }
}
