import { createFeatureSelector, createSelector } from '@ngrx/store';
import {
  RESORTS_FEATURE_KEY,
  ResortsState,
  resortsAdapter,
} from './resorts.reducer';

// Lookup the 'Resorts' feature state managed by NgRx
export const getResortsState =
  createFeatureSelector<ResortsState>(RESORTS_FEATURE_KEY);

const { selectAll, selectEntities } = resortsAdapter.getSelectors();

export const getResortsLoaded = createSelector(
  getResortsState,
  (state: ResortsState) => state.loaded
);

export const getResortsError = createSelector(
  getResortsState,
  (state: ResortsState) => state.error
);

export const getAllResorts = createSelector(
  getResortsState,
  (state: ResortsState) => selectAll(state)
);

export const getResortsEntities = createSelector(
  getResortsState,
  (state: ResortsState) => selectEntities(state)
);

export const getSelectedId = createSelector(
  getResortsState,
  (state: ResortsState) => state.selectedId
);

export const getSelected = createSelector(
  getResortsEntities,
  getSelectedId,
  (entities, selectedId) => (selectedId ? entities[selectedId] : undefined)
);
