import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Subject } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class RemoteService {
  clientID = Math.random().toString(36).substring(7);

  private protocol = window.location.protocol === 'https:' ? 'wss' : 'ws';
  ws = new WebSocket(
    `${this.protocol}://${window.location.host}/api/remote/ws/sub?clientID=${this.clientID}`
  );

  private _events = new Subject<Message<PossibleMessages>>();

  constructor(private httpClient: HttpClient) {
    this.ws.onmessage = (event) => {
      this.processEvent(event);
    };

    this.ws.onopen = () => {
      console.log('Connected');
    };

    this.ws.onclose = () => {
      console.log('Disconnected');
    };

    this.ws.onerror = (err) => {
      console.error(err);
    };
  }

  post(msg: PossibleMessages) {
    this.httpClient
      .post(`/api/remote/pub?message=${msg}&clientID=${this.clientID}`, null)
      .subscribe();
  }

  events() {
    return this._events.asObservable();
  }

  private processEvent(evt: MessageEvent) {
    let data: Message<PossibleMessages>;
    try {
      const rawMsg = JSON.parse(evt.data) as Message<string>;
      data = {
        UUID: rawMsg.UUID,
        Payload: atob(rawMsg.Payload) as PossibleMessages,
        Metadata: rawMsg.Metadata,
      };
    } catch (error) {
      console.error('Failed to parse message', error);
      return;
    }

    this._events.next(data);
  }
}

interface Message<T> {
  UUID: string;
  Payload: T;
  Metadata: Record<string, any>;
}

export type PossibleMessages =
  | 'playpause'
  | 'next'
  | 'previous'
  | 'revealhide'
  | 'refresh';
