import { Routes } from '@angular/router';
import { HomeComponent } from './home/home.component';
import { RemoteControlComponent } from './remote-control/remote-control.component';

export const routes: Routes = [
  {
    path: '',
    component: HomeComponent,
  },
  {
    path: 'remote',
    component: RemoteControlComponent,
  },
];
