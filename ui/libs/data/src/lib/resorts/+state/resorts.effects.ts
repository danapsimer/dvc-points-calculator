import { Injectable } from '@angular/core';
import { createEffect, Actions, ofType } from '@ngrx/effects';
import { fetch } from '@nrwl/angular';

import * as ResortsActions from './resorts.actions';
import * as ResortsFeature from './resorts.reducer';

@Injectable()
export class ResortsEffects {
  init$ = createEffect(() =>
    this.actions$.pipe(
      ofType(ResortsActions.initResorts),
      fetch({
        run: (action) => {
          // Your custom service 'load' logic goes here. For now just return a success action...
          return ResortsActions.loadResortsSuccess({ resorts: [] });
        },
        onError: (action, error) => {
          console.error('Error', error);
          return ResortsActions.loadResortsFailure({ error });
        },
      })
    )
  );

  constructor(private readonly actions$: Actions) {}
}
