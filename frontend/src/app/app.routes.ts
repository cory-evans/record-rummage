import { Routes } from '@angular/router';
import { HomeComponent } from './home/home.component';
import { RemoteControlComponent } from './remote-control/remote-control.component';
import { SettingsComponent } from './settings/settings.component';
import { authGuard } from './shared/guards/auth.guard';

export const routes: Routes = [
  {
    path: 'remote',
    component: RemoteControlComponent,
    canActivate: [authGuard],
  },
  {
    path: 'settings',
    component: SettingsComponent,
    canActivate: [authGuard],
  },
  {
    path: 'library',
    canActivate: [authGuard],
    loadComponent: () =>
      import('./library/library.component').then((m) => m.LibraryComponent),
  },
  {
    path: '',
    component: HomeComponent,
  },
];
