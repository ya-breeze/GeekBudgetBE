import { Injectable, signal } from '@angular/core';

@Injectable({
  providedIn: 'root',
})
export class LayoutService {
  readonly sidenavOpened = signal(true);
  readonly sidenavWidth = 250; // Width in pixels

  toggleSidenav(): void {
    this.sidenavOpened.update((value) => !value);
  }
}

