import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';
import { AuthService } from '../services/auth.service';
import { tap } from 'rxjs';

export const authGuard: CanActivateFn = (route, state) => {
  const authService = inject(AuthService);
  const router = inject(Router);

  return authService.signedIn$.pipe(
    tap((isSignedIn) => {
      if (!isSignedIn) {
        router.navigate(['/']);
      }
    })
  );
};
