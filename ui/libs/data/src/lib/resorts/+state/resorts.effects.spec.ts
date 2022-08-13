import { TestBed } from '@angular/core/testing';
import { provideMockActions } from '@ngrx/effects/testing';
import { Action } from '@ngrx/store';
import { provideMockStore } from '@ngrx/store/testing';
import { NxModule } from '@nrwl/angular';
import { hot } from 'jasmine-marbles';
import { Observable } from 'rxjs';

import * as ResortsActions from './resorts.actions';
import { ResortsEffects } from './resorts.effects';

describe('ResortsEffects', () => {
  let actions: Observable<Action>;
  let effects: ResortsEffects;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [NxModule.forRoot()],
      providers: [
        ResortsEffects,
        provideMockActions(() => actions),
        provideMockStore(),
      ],
    });

    effects = TestBed.inject(ResortsEffects);
  });

  describe('init$', () => {
    it('should work', () => {
      actions = hot('-a-|', { a: ResortsActions.initResorts() });

      const expected = hot('-a-|', {
        a: ResortsActions.loadResortsSuccess({ resorts: [] }),
      });

      expect(effects.init$).toBeObservable(expected);
    });
  });
});
