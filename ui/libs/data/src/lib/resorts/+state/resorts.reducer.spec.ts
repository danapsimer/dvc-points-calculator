import { Action } from '@ngrx/store';

import * as ResortsActions from './resorts.actions';
import { ResortsEntity } from './resorts.models';
import {
  ResortsState,
  initialResortsState,
  resortsReducer,
} from './resorts.reducer';

describe('Resorts Reducer', () => {
  const createResortsEntity = (id: string, name = ''): ResortsEntity => ({
    id,
    name: name || `name-${id}`,
  });

  describe('valid Resorts actions', () => {
    it('loadResortsSuccess should return the list of known Resorts', () => {
      const resorts = [
        createResortsEntity('PRODUCT-AAA'),
        createResortsEntity('PRODUCT-zzz'),
      ];
      const action = ResortsActions.loadResortsSuccess({ resorts });

      const result: ResortsState = resortsReducer(initialResortsState, action);

      expect(result.loaded).toBe(true);
      expect(result.ids.length).toBe(2);
    });
  });

  describe('unknown action', () => {
    it('should return the previous state', () => {
      const action = {} as Action;

      const result = resortsReducer(initialResortsState, action);

      expect(result).toBe(initialResortsState);
    });
  });
});
