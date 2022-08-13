import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { StoreModule } from '@ngrx/store';
import { EffectsModule } from '@ngrx/effects';
import * as fromResorts from './+state/resorts.reducer';
import { ResortsEffects } from './+state/resorts.effects';
import { ResortsFacade } from './+state/resorts.facade';

@NgModule({
  declarations: [],
  imports: [
    CommonModule,
    StoreModule.forFeature(
      fromResorts.RESORTS_FEATURE_KEY,
      fromResorts.resortsReducer
    ),
    EffectsModule.forFeature([ResortsEffects]),
  ],
  providers: [ResortsFacade],
})
export class ResortsModule {}
