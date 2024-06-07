import { AsyncPipe, JsonPipe } from '@angular/common';
import { HttpClient, HttpClientModule } from '@angular/common/http';
import { TmplAstSwitchBlockCase } from '@angular/compiler';
import { Component } from '@angular/core';
import { map, shareReplay } from 'rxjs';
import { NowPlayingComponent } from '../now-playing/now-playing.component';
import { SharedModule } from '../shared/shared.module';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [
    HttpClientModule,
    AsyncPipe,
    JsonPipe,
    NowPlayingComponent,
    SharedModule,
  ],
  templateUrl: './home.component.html',
})
export class HomeComponent {
  constructor(private readonly http: HttpClient) {}

  private baseUrl = '/api/';

  login() {
    window.location.href = this.baseUrl + 'auth/login';
  }

  logout() {
    window.location.href = this.baseUrl + 'auth/logout';
  }

  me$ = this.http.get(this.baseUrl + 'auth/me').pipe(shareReplay(1));

  signedIn$ = this.me$.pipe(map((me) => !!me));
}
