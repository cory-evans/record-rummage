import { Component } from '@angular/core';
import { SettingsService } from '../settings/settings.service';
import { SharedModule } from '../shared/shared.module';

@Component({
  selector: 'app-unauthenticated',
  standalone: true,
  imports: [SharedModule],
  templateUrl: './unauthenticated.component.html',
})
export class UnauthenticatedComponent {
  constructor(public readonly settingsService: SettingsService) {}
  login() {
    window.location.href = this.settingsService.baseAPIUrl + 'auth/login';
  }
}
