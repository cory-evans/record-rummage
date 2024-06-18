import { CommonModule } from '@angular/common';
import { HttpClient, HttpClientModule } from '@angular/common/http';
import { Component } from '@angular/core';
import {
  BehaviorSubject,
  Observable,
  Subject,
  combineLatest,
  first,
  firstValueFrom,
  map,
  of,
  shareReplay,
  switchMap,
  tap,
  timer,
} from 'rxjs';
import { SimplePlaylistPage } from '../../types/spotify';
import { SharedModule } from '../shared/shared.module';

@Component({
  selector: 'app-library',
  standalone: true,
  imports: [CommonModule, HttpClientModule, SharedModule],
  templateUrl: './library.component.html',
})
export class LibraryComponent {
  page$ = new BehaviorSubject<number>(1);

  playlists$ = this.page$.pipe(
    switchMap(
      (n) => this.http.get<SimplePlaylistPage>(`/api/playlist/mine?page=${n}`)
      // of(DATA)
    ),
    shareReplay(1)
  );

  items$ = this.playlists$.pipe(map((page) => page.items));

  nPages$ = combineLatest([this.playlists$, this.page$]).pipe(
    map(([p, current]) => {
      const items = Array.from(
        { length: Math.ceil(p.total / p.limit) },
        (_, i) => i + 1
      );

      return items.map((i) => ({ idx: i, active: i === current }));
    })
  );

  constructor(private readonly http: HttpClient) {}

  loading = false;
  progress$ = new BehaviorSubject<number>(0);

  async refresh(playlistID: string) {
    this.loading = true;
    await firstValueFrom(
      of(null).pipe(
        tap(() => this.progress$.next(0)),
        switchMap(() =>
          this.http.post('/api/playlist/refresh?id=' + playlistID, null)
        ),
        switchMap(() => timer(0, 250)),
        switchMap(() =>
          this.http.get<{ progress: number }>(
            '/api/playlist/refresh-progress?id=' + playlistID
          )
        ),
        tap((p) => {
          this.progress$.next(p.progress);
        }),
        first((p) => p.progress === 1)
      )
    );

    this.loading = false;
  }
}
