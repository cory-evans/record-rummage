import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { SettingsService } from './settings/settings.service';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet],
  templateUrl: './app.component.html',
  styleUrl: './app.component.css',
  host: {
    class: 'flex-1 flex flex-col',
  },
})
export class AppComponent {
  title = 'frontend';

  constructor(settings: SettingsService) {}
}
