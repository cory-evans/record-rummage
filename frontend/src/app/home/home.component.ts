import { CommonModule } from '@angular/common';
import { HttpClient } from '@angular/common/http';
import { Component } from '@angular/core';
import { NowPlayingComponent } from '../now-playing/now-playing.component';
import { SharedModule } from '../shared/shared.module';
import { RouterModule } from '@angular/router';
import { SettingsService } from '../settings/settings.service';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [CommonModule, SharedModule, RouterModule, NowPlayingComponent],
  templateUrl: './home.component.html',
  host: {
    class: 'flex-1 flex flex-col',
  },
})
export class HomeComponent {
  constructor(
    private readonly http: HttpClient,
    public readonly settings: SettingsService
  ) {}
}
