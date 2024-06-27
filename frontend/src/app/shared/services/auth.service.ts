import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { SettingsService } from '../../settings/settings.service';
import { catchError, map, of, shareReplay } from 'rxjs';
import { Router } from '@angular/router';

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  constructor(
    private readonly http: HttpClient,
    private readonly settings: SettingsService,
    private readonly router: Router
  ) {}

  me$ = this.http.get(this.settings.baseAPIUrl + 'auth/me').pipe(
    catchError(() => {
      return of(null);
    }),
    shareReplay(1)
  );

  signedIn$ = this.me$.pipe(map((me) => !!me));
}
