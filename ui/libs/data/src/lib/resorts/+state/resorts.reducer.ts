import { EntityState, EntityAdapter, createEntityAdapter } from '@ngrx/entity';
import { createReducer, on, Action } from '@ngrx/store';

import * as ResortsActions from './resorts.actions';
import { ResortsEntity } from './resorts.models';

export const RESORTS_FEATURE_KEY = 'resorts';

export interface ResortsState extends EntityState<ResortsEntity> {
  selectedId?: string | number; // which Resorts record has been selected
  loaded: boolean; // has the Resorts list been loaded
  error?: string | null; // last known error (if any)
}

export interface ResortsPartialState {
  readonly [RESORTS_FEATURE_KEY]: ResortsState;
}

export const resortsAdapter: EntityAdapter<ResortsEntity> =
  createEntityAdapter<ResortsEntity>();

export const initialResortsState: ResortsState = resortsAdapter.getInitialState(
  {
    // set initial required properties
    loaded: false,
  }
);

const reducer = createReducer(
  initialResortsState,
  on(ResortsActions.initResorts, (state) => ({
    ...state,
    loaded: false,
    error: null,
  })),
  on(ResortsActions.loadResortsSuccess, (state, { resorts }) =>
    resortsAdapter.setAll(resorts, { ...state, loaded: true })
  ),
  on(ResortsActions.loadResortsFailure, (state, { error }) => ({
    ...state,
    error,
  }))
);

export function resortsReducer(
  state: ResortsState | undefined,
  action: Action
) {
  return reducer(state, action);
}
