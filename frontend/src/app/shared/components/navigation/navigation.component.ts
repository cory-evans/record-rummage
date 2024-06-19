import { Component, ElementRef, Inject, viewChild } from '@angular/core';
import { SettingsService } from '../../../settings/settings.service';
import { NowPlayingService } from '../../../now-playing/now-playing.service';
import { DOCUMENT } from '@angular/common';

@Component({
  selector: 'app-navigation',
  templateUrl: './navigation.component.html',
})
export class NavigationComponent {
  menuBtn = viewChild<ElementRef>('menuBtn');

  constructor(
    public readonly settings: SettingsService,
    public readonly nowPlaying: NowPlayingService,
    @Inject(DOCUMENT) private document: Document
  ) {}

  login() {
    window.location.href = this.settings.baseAPIUrl + 'auth/login';
  }

  logout() {
    window.location.href = this.settings.baseAPIUrl + 'auth/logout';
  }

  blur() {
    if (this.document.activeElement instanceof HTMLElement) {
      this.document.activeElement.blur();
    }
  }
}
