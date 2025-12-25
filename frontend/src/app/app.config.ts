import {
    ApplicationConfig,
    provideBrowserGlobalErrorListeners,
    provideZoneChangeDetection,
    APP_INITIALIZER,
    inject,
} from '@angular/core';
import { provideRouter } from '@angular/router';
import { provideHttpClient, withInterceptors } from '@angular/common/http';
import { provideAnimationsAsync } from '@angular/platform-browser/animations/async';
import { provideCharts, withDefaultRegisterables } from 'ng2-charts';
import { provideNativeDateAdapter } from '@angular/material/core';

import { routes } from './app.routes';
import { AuthService } from './core/auth/services/auth.service';
// import { authInterceptor } from './core/auth/interceptors/auth.interceptor';
import { errorInterceptor } from './core/auth/interceptors/error.interceptor';
import { ApiConfiguration } from './core/api/api-configuration';
import { environment } from '../environments/environment';

export const appConfig: ApplicationConfig = {
    providers: [
        provideBrowserGlobalErrorListeners(),
        provideZoneChangeDetection({ eventCoalescing: true }),
        provideRouter(routes),
        provideHttpClient(withInterceptors([errorInterceptor])),
        provideAnimationsAsync(),
        provideCharts(withDefaultRegisterables()),
        provideNativeDateAdapter(),
        {
            provide: APP_INITIALIZER,
            useFactory: () => {
                const authService = inject(AuthService);
                return () => authService.checkAuth();
            },
            multi: true,
        },
        {
            provide: ApiConfiguration,
            useFactory: () => {
                const config = new ApiConfiguration();
                config.rootUrl = environment.apiUrl;
                return config;
            },
        },
    ],
};
