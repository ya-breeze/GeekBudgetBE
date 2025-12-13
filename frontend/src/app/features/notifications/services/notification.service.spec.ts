import { TestBed } from '@angular/core/testing';
import { HttpClientTestingModule, HttpTestingController } from '@angular/common/http/testing';
import { NotificationService } from './notification.service';
import { ApiConfiguration } from '../../../core/api/api-configuration';
import { Notification } from '../../../core/api/models/notification';

describe('NotificationService', () => {
  let service: NotificationService;
  let httpMock: HttpTestingController;
  let apiConfig: ApiConfiguration;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule],
      providers: [NotificationService, ApiConfiguration],
    });

    service = TestBed.inject(NotificationService);
    httpMock = TestBed.inject(HttpTestingController);
    apiConfig = TestBed.inject(ApiConfiguration);
  });

  afterEach(() => {
    httpMock.verify();
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  describe('loadNotifications', () => {
    it('should load all notifications', (done) => {
      const mockNotifications: Notification[] = [
        {
          id: '1',
          date: '2024-01-01',
          type: 'balanceMatch',
          title: 'Balance Match',
          description: 'Balance matches expected value',
        },
      ];

      service.loadNotifications().subscribe({
        next: (notifications) => {
          expect(notifications).toEqual(mockNotifications);
          expect(service.notifications()).toEqual(mockNotifications);
          done();
        },
      });

      const req = httpMock.expectOne(`${apiConfig.rootUrl}/v1/notifications`);
      expect(req.request.method).toBe('GET');
      req.flush(mockNotifications);
    });

    it('should handle different notification types', (done) => {
      const mockNotifications: Notification[] = [
        {
          id: '1',
          date: '2024-01-01',
          type: 'balanceMatch',
          title: 'Balance Match',
          description: 'Balance matches',
        },
        {
          id: '2',
          date: '2024-01-02',
          type: 'balanceDoesntMatch',
          title: 'Balance Mismatch',
          description: 'Balance does not match',
        },
        {
          id: '3',
          date: '2024-01-03',
          type: 'other',
          title: 'Other Notification',
          description: 'General notification',
        },
      ];

      service.loadNotifications().subscribe({
        next: (notifications) => {
          expect(notifications.length).toBe(3);
          expect(notifications[0].type).toBe('balanceMatch');
          expect(notifications[1].type).toBe('balanceDoesntMatch');
          expect(notifications[2].type).toBe('other');
          done();
        },
      });

      const req = httpMock.expectOne(`${apiConfig.rootUrl}/v1/notifications`);
      req.flush(mockNotifications);
    });

    it('should update signal state correctly', (done) => {
      const mockNotifications: Notification[] = [
        {
          id: '1',
          date: '2024-01-01',
          type: 'other',
          title: 'Test',
          description: 'Test notification',
        },
      ];

      service.loadNotifications().subscribe({
        next: () => {
          expect(service.notifications()).toEqual(mockNotifications);
          expect(service.loading()).toBe(false);
          expect(service.error()).toBeNull();
          done();
        },
      });

      const req = httpMock.expectOne(`${apiConfig.rootUrl}/v1/notifications`);
      req.flush(mockNotifications);
    });
  });

  describe('delete', () => {
    it('should delete a notification', (done) => {
      service.delete('1').subscribe({
        next: () => {
          done();
        },
      });

      const req = httpMock.expectOne(`${apiConfig.rootUrl}/v1/notifications/1`);
      expect(req.request.method).toBe('DELETE');
      req.flush(null);
    });
  });
});
