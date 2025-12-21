import { TestBed } from '@angular/core/testing';
import { HttpClientTestingModule, HttpTestingController } from '@angular/common/http/testing';
import { AuthService } from './auth.service';
import { ApiConfiguration } from '../../api/api-configuration';

describe('AuthService', () => {
    let service: AuthService;
    let httpMock: HttpTestingController;
    let apiConfig: ApiConfiguration;

    beforeEach(() => {
        TestBed.configureTestingModule({
            imports: [HttpClientTestingModule],
            providers: [AuthService, ApiConfiguration],
        });

        service = TestBed.inject(AuthService);
        httpMock = TestBed.inject(HttpTestingController);
        apiConfig = TestBed.inject(ApiConfiguration);
    });

    afterEach(() => {
        httpMock.verify();
        localStorage.clear();
    });

    it('should be created', () => {
        expect(service).toBeTruthy();
    });

    describe('login', () => {
        it('should login successfully with valid credentials', (done) => {
            const mockCredentials = { email: 'test@example.com', password: 'password123' };
            const mockResponse = { token: 'mock-jwt-token' }; // Backend still returns token, but service ignores it

            service.login(mockCredentials.email, mockCredentials.password).subscribe({
                next: () => {
                    expect(service.isLoggedIn()).toBe(true);
                    done();
                },
            });

            const req = httpMock.expectOne(`${apiConfig.rootUrl}/v1/authorize`);
            expect(req.request.method).toBe('POST');
            expect(req.request.body).toEqual(mockCredentials);
            req.flush(mockResponse);
        });

        it('should update authentication state after successful login', (done) => {
            const mockResponse = { token: 'mock-jwt-token' };

            service.login('test@example.com', 'password123').subscribe({
                next: () => {
                    expect(service.isLoggedIn()).toBe(true);
                    done();
                },
            });

            const req = httpMock.expectOne(`${apiConfig.rootUrl}/v1/authorize`);
            req.flush(mockResponse);
        });

        it('should return error on invalid credentials', (done) => {
            const mockError = { status: 401, statusText: 'Unauthorized' };

            service.login('test@example.com', 'wrong-password').subscribe({
                error: (error) => {
                    expect(error.status).toBe(401);
                    done();
                },
            });

            const req = httpMock.expectOne(`${apiConfig.rootUrl}/v1/authorize`);
            req.flush('Unauthorized', mockError);
        });
    });

    describe('logout', () => {
        it('should logout and clear state', () => {
            // Simulate logged in state
            // service.isAuthenticatedSubject.next(true); // Private access workaround or just login
            // We can't access private subject easily, but we can call login or rely on public api if we could set state.
            // For now, let's assume it works if we verify logout actions.

            // Trigger login to set state to true
            const mockResponse = { token: 'mock-token' };
            service.login('a', 'b').subscribe();
            const req = httpMock.expectOne(`${apiConfig.rootUrl}/v1/authorize`);
            req.flush(mockResponse);
            expect(service.isLoggedIn()).toBe(true);

            service.logout();
            expect(service.isLoggedIn()).toBe(false);

            const reqLogout = httpMock.expectOne(`${apiConfig.rootUrl}/v1/logout`);
            expect(reqLogout.request.method).toBe('POST');
            reqLogout.flush({});
        });
    });

    describe('isLoggedIn', () => {
        it('should return false initially', () => {
            // Re-inject to ensure fresh state
            service = TestBed.inject(AuthService);
            expect(service.isLoggedIn()).toBe(false);
        });
    });

    describe('checkAuth', () => {
        it('should return true and set authenticated state on success', (done) => {
            service.checkAuth().subscribe((isLoggedIn) => {
                expect(isLoggedIn).toBe(true);
                expect(service.isLoggedIn()).toBe(true);
                done();
            });

            const req = httpMock.expectOne(`${apiConfig.rootUrl}/v1/user`);
            expect(req.request.method).toBe('GET');
            req.flush({}); // Return empty user object or whatever, success 200
        });

        it('should return false and set unauthenticated state on failure', (done) => {
            service.checkAuth().subscribe((isLoggedIn) => {
                expect(isLoggedIn).toBe(false);
                expect(service.isLoggedIn()).toBe(false);
                done();
            });

            const req = httpMock.expectOne(`${apiConfig.rootUrl}/v1/user`);
            req.flush('Unauthorized', { status: 401, statusText: 'Unauthorized' });
        });
    });
});
