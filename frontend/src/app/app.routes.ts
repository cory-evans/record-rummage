import { Routes } from '@angular/router';
import { HomeComponent } from './home/home.component';
import { RemoteControlComponent } from './remote-control/remote-control.component';
import { SettingsComponent } from './settings/settings.component';

export const routes: Routes = [
  {
    path: 'remote',
    component: RemoteControlComponent,
  },
  {
    path: 'settings',
    component: SettingsComponent,
  },
  {
    path: 'library',
    loadComponent: () =>
      import('./library/library.component').then((m) => m.LibraryComponent),
  },
  {
    path: '',
    component: HomeComponent,
  },
];
