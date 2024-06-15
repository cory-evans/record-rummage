import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { IconComponent } from './components/icon/icon.component';
import { RouterModule } from '@angular/router';

@NgModule({
  declarations: [IconComponent],
  imports: [CommonModule, RouterModule],
  exports: [IconComponent],
})
export class SharedModule {}
