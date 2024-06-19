import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { SettingsService } from './settings/settings.service';
import { SharedModule } from './shared/shared.module';
import { AuthService } from './shared/services/auth.service';
import { AsyncPipe } from '@angular/common';
import { UnauthenticatedComponent } from './unauthenticated/unauthenticated.component';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet, SharedModule, AsyncPipe, UnauthenticatedComponent],
  templateUrl: './app.component.html',
  styleUrl: './app.component.css',
  host: {
    class: 'flex-1 flex flex-col',
  },
})
export class AppComponent {
  title = 'frontend';

  constructor(public auth: AuthService) {}
}
