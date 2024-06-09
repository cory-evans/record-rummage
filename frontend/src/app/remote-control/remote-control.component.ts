import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { SharedModule } from '../shared/shared.module';
import { RouterModule } from '@angular/router';
import {
  PossibleMessages,
  RemoteService,
} from '../shared/services/remote.service';

@Component({
  selector: 'app-remote-control',
  standalone: true,
  imports: [CommonModule, SharedModule, RouterModule],
  templateUrl: './remote-control.component.html',
  host: {
    class: 'flex-1 flex flex-col',
  },
})
export class RemoteControlComponent {
  constructor(private readonly remoteService: RemoteService) {
    remoteService.events().subscribe((evt) => {
      console.log('Event:', evt);
    });
  }

  send(msg: PossibleMessages) {
    this.remoteService.post(msg);
  }
}
