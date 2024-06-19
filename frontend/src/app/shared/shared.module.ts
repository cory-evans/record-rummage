import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { IconComponent } from './components/icon/icon.component';
import { RouterModule } from '@angular/router';
import { NavigationComponent } from './components/navigation/navigation.component';

@NgModule({
  declarations: [IconComponent, NavigationComponent],
  imports: [CommonModule, RouterModule],
  exports: [IconComponent, NavigationComponent],
})
export class SharedModule {}
