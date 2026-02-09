import { CommonModule } from '@angular/common';
import { HttpClientModule } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { environment } from '../../../environments/environment';
import { BackendStatus, StatusService } from '../../core/services/status.service';

@Component({
    selector: 'app-footer',
    standalone: true,
    imports: [CommonModule, HttpClientModule],
    templateUrl: './footer.component.html',
    styleUrl: './footer.component.scss',
})
export class FooterComponent implements OnInit {
    backendStatus?: BackendStatus;
    frontendBuildTime = environment.buildTime;

    constructor(private statusService: StatusService) {}

    ngOnInit(): void {
        this.statusService.getStatus().subscribe({
            next: (status) => (this.backendStatus = status),
            error: (err) => console.error('Failed to fetch backend status', err),
        });
    }
}
