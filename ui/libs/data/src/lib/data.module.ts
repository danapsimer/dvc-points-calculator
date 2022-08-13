import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ResortsModule } from './resorts/resorts.module';
import { StoreModule } from '@ngrx/store';
import { EffectsModule } from '@ngrx/effects';
import { StoreDevtoolsModule } from '@ngrx/store-devtools';
import { environment } from '../../../../apps/points/src/environments/environment';

@NgModule({
  imports: [
    CommonModule,
    ResortsModule,
    StoreModule.forRoot(
      {},
      {
        metaReducers: !environment.production ? [] : [],
        runtimeChecks: {
          strictActionImmutability: true,
          strictStateImmutability: true,
        },
      }
    ),
    EffectsModule.forRoot([]),
    !environment.production ? StoreDevtoolsModule.instrument() : [],
  ],
})
export class DataModule {}
