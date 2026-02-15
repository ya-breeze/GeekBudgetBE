import { Component, inject, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { AuthService } from '../auth/services/auth.service';

@Component({
    selector: 'app-redirect',
    standalone: true,
    template: '',
})
export class RedirectComponent implements OnInit {
    private readonly authService = inject(AuthService);
    private readonly router = inject(Router);

    ngOnInit(): void {
        // Redirect based on authentication status
        if (this.authService.isLoggedIn()) {
            this.router.navigate(['/dashboard']);
        } else {
            this.router.navigate(['/landing']);
        }
    }
}
