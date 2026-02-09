import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';

export interface BackendStatus {
    buildTime: string;
    commit: string;
    startTime: string;
}

@Injectable({
    providedIn: 'root',
})
export class StatusService {
    constructor(private http: HttpClient) {}

    getStatus(): Observable<BackendStatus> {
        const url = `${environment.apiUrl}/${environment.apiVersion}/status`;
        return this.http.get<BackendStatus>(url);
    }
}
