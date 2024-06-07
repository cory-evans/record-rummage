import { CommonModule } from '@angular/common';
import { Component, PLATFORM_ID } from '@angular/core';
import {
  BehaviorSubject,
  combineLatest,
  firstValueFrom,
  of,
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

  constructor(private readonly http: HttpClient) {}

  currentlyPlaying$ = combineLatest([this.timer$, this.forceRefresh$]).pipe(
    switchMap(() => this.http.get<Root>('/api/track/currently-playing')),
    shareReplay(1)
  );

  refreshPlaylist() {
    // let id = '4CGbyFVJTDvQDDa7Pg8IaO';
    // this.http.post('/api/playlist/refresh?id=' + id, {}).subscribe(() => {});
  }

  async togglePlayPause(is_playing: boolean) {
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
  user: {
    id: string;
    display_name: string;
    images: {
      url: string;
      height: number;
      width: number;
    }[];
  } | null = null;
  async reveal() {
    const track = await firstValueFrom(this.currentlyPlaying$);

    await firstValueFrom(
      this.http
        .get<typeof this.user>(
          `/api/track/reveal?trackId=${track.item.id}&playlistId=${this.PlaylistID}`
        )
        .pipe(
          tap((resp) => {
            if (!resp) {
              console.log('Track not found in playlist');
              this.user = null;
              return;
            }

            this.lastRevealedTrackId = track.item.id;

            this.user = resp;
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
