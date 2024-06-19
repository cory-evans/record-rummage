import { DOCUMENT } from '@angular/common';
import { Inject, Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root',
})
export class SettingsService {
  constructor(@Inject(DOCUMENT) private readonly document: Document) {
    this.updateDOM();
  }

  get theme(): ThemeSetting {
    const x = localStorage.getItem('theme');
    return x ? (x as ThemeSetting) : 'cupcake';
  }

  setTheme(theme: ThemeSetting) {
    localStorage.setItem('theme', theme);

    this.updateDOM();
  }

  toggleTheme(event: Event) {
    const t = event.target as HTMLInputElement;
    const theme = t?.checked ? 'dark' : 'light';
    this.setTheme(theme);
  }

  private updateDOM() {
    const html = this.document.querySelector('html');
    if (html) {
      html.setAttribute('data-theme', this.theme);
    }
  }
}

type ThemeSetting = 'light' | 'dark' | 'cupcake';
