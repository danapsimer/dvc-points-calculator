import { NgModule } from '@angular/core';
import { TestBed } from '@angular/core/testing';
import { EffectsModule } from '@ngrx/effects';
import { StoreModule, Store } from '@ngrx/store';
import { NxModule } from '@nrwl/angular';
import { readFirst } from '@nrwl/angular/testing';

import * as ResortsActions from './resorts.actions';
import { ResortsEffects } from './resorts.effects';
import { ResortsFacade } from './resorts.facade';
import { ResortsEntity } from './resorts.models';
import {
  RESORTS_FEATURE_KEY,
  ResortsState,
  initialResortsState,
  resortsReducer,
} from './resorts.reducer';
import * as ResortsSelectors from './resorts.selectors';

interface TestSchema {
  resorts: ResortsState;
}

describe('ResortsFacade', () => {
  let facade: ResortsFacade;
  let store: Store<TestSchema>;
  const createResortsEntity = (id: string, name = ''): ResortsEntity => ({
    id,
    name: name || `name-${id}`,
  });

  describe('used in NgModule', () => {
    beforeEach(() => {
      @NgModule({
        imports: [
          StoreModule.forFeature(RESORTS_FEATURE_KEY, resortsReducer),
          EffectsModule.forFeature([ResortsEffects]),
        ],
        providers: [ResortsFacade],
      })
      class CustomFeatureModule {}

      @NgModule({
        imports: [
          NxModule.forRoot(),
          StoreModule.forRoot({}),
          EffectsModule.forRoot([]),
          CustomFeatureModule,
        ],
      })
      class RootModule {}
      TestBed.configureTestingModule({ imports: [RootModule] });

      store = TestBed.inject(Store);
      facade = TestBed.inject(ResortsFacade);
    });

    /**
     * The initially generated facade::loadAll() returns empty array
     */
    it('loadAll() should return empty list with loaded == true', async () => {
      let list = await readFirst(facade.allResorts$);
      let isLoaded = await readFirst(facade.loaded$);

      expect(list.length).toBe(0);
      expect(isLoaded).toBe(false);

      facade.init();

      list = await readFirst(facade.allResorts$);
      isLoaded = await readFirst(facade.loaded$);

      expect(list.length).toBe(0);
      expect(isLoaded).toBe(true);
    });

    /**
     * Use `loadResortsSuccess` to manually update list
     */
    it('allResorts$ should return the loaded list; and loaded flag == true', async () => {
      let list = await readFirst(facade.allResorts$);
      let isLoaded = await readFirst(facade.loaded$);

      expect(list.length).toBe(0);
      expect(isLoaded).toBe(false);

      store.dispatch(
        ResortsActions.loadResortsSuccess({
          resorts: [createResortsEntity('AAA'), createResortsEntity('BBB')],
        })
      );

      list = await readFirst(facade.allResorts$);
      isLoaded = await readFirst(facade.loaded$);

      expect(list.length).toBe(2);
      expect(isLoaded).toBe(true);
    });
  });
});
