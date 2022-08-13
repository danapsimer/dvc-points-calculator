import { Injectable } from '@angular/core';
import { select, Store, Action } from '@ngrx/store';

import * as ResortsActions from './resorts.actions';
import * as ResortsFeature from './resorts.reducer';
import * as ResortsSelectors from './resorts.selectors';

@Injectable()
export class ResortsFacade {
  /**
   * Combine pieces of state using createSelector,
   * and expose them as observables through the facade.
   */
  loaded$ = this.store.pipe(select(ResortsSelectors.getResortsLoaded));
  allResorts$ = this.store.pipe(select(ResortsSelectors.getAllResorts));
  selectedResorts$ = this.store.pipe(select(ResortsSelectors.getSelected));

  constructor(private readonly store: Store) {}

  /**
   * Use the initialization action to perform one
   * or more tasks in your Effects.
   */
  init() {
    this.store.dispatch(ResortsActions.initResorts());
  }
}
