import { Injectable } from '@angular/core';
import { Subject } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class NowPlayingService {
  private _refresh$ = new Subject<void>();

  constructor() {}

  events() {
    return this._refresh$.asObservable();
  }

  refresh() {
    this._refresh$.next();
  }
}
