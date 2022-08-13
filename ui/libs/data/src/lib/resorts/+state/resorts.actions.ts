import { createAction, props } from '@ngrx/store';
import { ResortsEntity } from './resorts.models';

export const initResorts = createAction('[Resorts Page] Init');

export const loadResortsSuccess = createAction(
  '[Resorts/API] Load Resorts Success',
  props<{ resorts: ResortsEntity[] }>()
);

export const loadResortsFailure = createAction(
  '[Resorts/API] Load Resorts Failure',
  props<{ error: any }>()
);
