import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import {
  BehaviorSubject,
  Subject,
  combineLatest,
  filter,
  firstValueFrom,
  map,
  shareReplay,
  switchMap,
  tap,
  timer,
} from 'rxjs';
import { SharedModule } from '../shared/shared.module';
import { HttpClient, HttpClientModule } from '@angular/common/http';

@Component({
  selector: 'app-now-playing',
  standalone: true,
  imports: [CommonModule, SharedModule, HttpClientModule],
  templateUrl: './now-playing.component.html',
})
export class NowPlayingComponent {
  readonly PlaylistID = '4CGbyFVJTDvQDDa7Pg8IaO';

  timer$ = timer(0, 60_000);
  forceRefresh$ = new BehaviorSubject<void>(void 0);

  constructor(private readonly http: HttpClient) {
    // every second check to see if we should force refresh
    timer(1_000, 1_000)
      .pipe(
        switchMap(() => this.dueToEnd$),
        map((dueToEnd) => {
          return dueToEnd < Date.now();
        }),
        filter((x) => x)
      )
      .subscribe(() => {
        this.forceRefresh$.next();
      });
  }

  private currentlyPlayingResponseTime$ = new BehaviorSubject<number>(
    Date.now()
  );
  currentlyPlaying$ = combineLatest([this.timer$, this.forceRefresh$]).pipe(
    switchMap(() => this.http.get<Root>('/api/track/currently-playing')),
    tap(() => this.currentlyPlayingResponseTime$.next(Date.now())),
    shareReplay(1)
  );

  dueToEnd$ = combineLatest([
    this.currentlyPlaying$,
    this.currentlyPlayingResponseTime$,
  ]).pipe(
    map(([x, t]) => {
      const msLeft = x.item.duration_ms - x.progress_ms;
      return t + msLeft;
    })
  );

  currentlyPlayingItem$ = this.currentlyPlaying$.pipe(map((x) => x.item));
  isPlaying$ = this.currentlyPlaying$.pipe(map((x) => x.is_playing));

  async togglePlayPause() {
    const is_playing = await firstValueFrom(this.isPlaying$);
    await firstValueFrom(
      this.http
        .put('/api/track/playback', { is_playing })
        .pipe(tap(() => this.forceRefresh$.next()))
    );
  }

  async nextTrack() {
    await firstValueFrom(
      this.http.post('/api/track/next', {}).pipe(
        tap(() => {
          setTimeout(() => {
            this.forceRefresh$.next();
          }, 1000);
        })
      )
    );
  }

  async previousTrack() {
    await firstValueFrom(
      this.http.post('/api/track/previous', {}).pipe(
        tap(() => {
          setTimeout(() => {
            this.forceRefresh$.next();
          }, 1000);
        })
      )
    );
  }

  lastRevealedTrackId: string | null = null;
  users:
    | {
        id: string;
        display_name: string;
        images: {
          url: string;
          height: number;
          width: number;
        }[];
      }[]
    | null = null;
  async reveal() {
    const track = await firstValueFrom(this.currentlyPlaying$);

    await firstValueFrom(
      this.http
        .get<typeof this.users>(
          `/api/track/reveal?trackId=${track.item.id}&playlistId=${this.PlaylistID}`
        )
        .pipe(
          tap((resp) => {
            if (!resp) {
              console.log('Track not found in playlist');
              this.users = null;
              return;
            }

            this.lastRevealedTrackId = track.item.id;

            this.users = resp;
          })
        )
    );
  }

  pickBigestImage(images: { url: string; height: number; width: number }[]) {
    return images.reduce((prev, curr) =>
      prev.height * prev.width > curr.height * curr.width ? prev : curr
    ).url;
  }
}
export interface Root {
  timestamp: number;
  context: Context;
  progress_ms: number;
  is_playing: boolean;
  item: Item;
}

export interface Context {
  external_urls: ExternalUrls;
  href: string;
  type: string;
  uri: string;
}

export interface ExternalUrls {
  spotify: string;
}

export interface Item {
  artists: Artist[];
  available_markets: string[];
  disc_number: number;
  duration_ms: number;
  explicit: boolean;
  external_urls: ExternalUrls3;
  href: string;
  id: string;
  name: string;
  preview_url: string;
  track_number: number;
  uri: string;
  type: string;
  album: Album;
  external_ids: ExternalIds;
  popularity: number;
  is_playable: any;
  linked_from: any;
}

export interface Artist {
  name: string;
  id: string;
  uri: string;
  href: string;
  external_urls: ExternalUrls2;
}

export interface ExternalUrls2 {
  spotify: string;
}

export interface ExternalUrls3 {
  spotify: string;
}

export interface Album {
  name: string;
  artists: Artist2[];
  album_group: string;
  album_type: string;
  id: string;
  uri: string;
  available_markets: string[];
  href: string;
  images: Image[];
  external_urls: ExternalUrls5;
  release_date: string;
  release_date_precision: string;
}

export interface Artist2 {
  name: string;
  id: string;
  uri: string;
  href: string;
  external_urls: ExternalUrls4;
}

export interface ExternalUrls4 {
  spotify: string;
}

export interface Image {
  height: number;
  width: number;
  url: string;
}

export interface ExternalUrls5 {
  spotify: string;
}

export interface ExternalIds {
  isrc: string;
}
