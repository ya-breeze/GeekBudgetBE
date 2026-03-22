import { Component, inject, input, output, OnInit, signal, computed } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatListModule } from '@angular/material/list';
import { MatInputModule } from '@angular/material/input';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { TemplateService } from '../services/template.service';
import { TransactionTemplate } from '../../../core/api/models/transaction-template';

@Component({
    selector: 'app-template-picker',
    standalone: true,
    imports: [
        FormsModule,
        MatListModule,
        MatInputModule,
        MatFormFieldModule,
        MatButtonModule,
        MatIconModule,
        MatProgressSpinnerModule,
    ],
    templateUrl: './template-picker.component.html',
})
export class TemplatePickerComponent implements OnInit {
    private readonly templateService = inject(TemplateService);

    /** Optional: pre-filter templates to those containing this accountId in their movements */
    accountId = input<string | undefined>(undefined);

    /** Emitted when the user selects a template */
    templateSelected = output<TransactionTemplate>();

    protected readonly loading = this.templateService.loading;
    protected readonly searchQuery = signal('');

    protected readonly filteredTemplates = computed(() => {
        const q = this.searchQuery().toLowerCase();
        return this.templateService.templates().filter((t) =>
            t.name.toLowerCase().includes(q),
        );
    });

    ngOnInit(): void {
        this.templateService.loadTemplates(this.accountId()).subscribe();
    }

    protected select(template: TransactionTemplate): void {
        this.templateSelected.emit(template);
    }
}
