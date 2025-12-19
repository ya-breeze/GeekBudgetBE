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
            const mockResponse = { token: 'mock-jwt-token' };

            service.login(mockCredentials.email, mockCredentials.password).subscribe({
                next: (token) => {
                    expect(token).toBe('mock-jwt-token');
                    done();
                },
            });

            const req = httpMock.expectOne(`${apiConfig.rootUrl}/v1/authorize`);
            expect(req.request.method).toBe('POST');
            expect(req.request.body).toEqual(mockCredentials);
            req.flush(mockResponse);
        });

        it('should store JWT token after successful login', (done) => {
            const mockResponse = { token: 'mock-jwt-token' };

            service.login('test@example.com', 'password123').subscribe({
                next: () => {
                    expect(localStorage.getItem('auth_token')).toBe('mock-jwt-token');
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
        it('should logout and clear token', () => {
            localStorage.setItem('auth_token', 'mock-token');
            service.logout();
            expect(localStorage.getItem('auth_token')).toBeNull();
        });
    });

    describe('isLoggedIn', () => {
        it('should return true when token exists', () => {
            localStorage.setItem('auth_token', 'mock-token');
            expect(service.isLoggedIn()).toBe(true);
        });

        it('should return false when token does not exist', () => {
            localStorage.removeItem('auth_token');
            expect(service.isLoggedIn()).toBe(false);
        });
    });

    describe('getToken', () => {
        it('should return token when it exists', () => {
            localStorage.setItem('auth_token', 'mock-token');
            expect(service.getToken()).toBe('mock-token');
        });

        it('should return null when token does not exist', () => {
            localStorage.removeItem('auth_token');
            expect(service.getToken()).toBeNull();
        });
    });
});
