import { ResortsEntity } from './resorts.models';
import {
  resortsAdapter,
  ResortsPartialState,
  initialResortsState,
} from './resorts.reducer';
import * as ResortsSelectors from './resorts.selectors';

describe('Resorts Selectors', () => {
  const ERROR_MSG = 'No Error Available';
  const getResortsId = (it: ResortsEntity) => it.id;
  const createResortsEntity = (id: string, name = '') =>
    ({
      id,
      name: name || `name-${id}`,
    } as ResortsEntity);

  let state: ResortsPartialState;

  beforeEach(() => {
    state = {
      resorts: resortsAdapter.setAll(
        [
          createResortsEntity('PRODUCT-AAA'),
          createResortsEntity('PRODUCT-BBB'),
          createResortsEntity('PRODUCT-CCC'),
        ],
        {
          ...initialResortsState,
          selectedId: 'PRODUCT-BBB',
          error: ERROR_MSG,
          loaded: true,
        }
      ),
    };
  });

  describe('Resorts Selectors', () => {
    it('getAllResorts() should return the list of Resorts', () => {
      const results = ResortsSelectors.getAllResorts(state);
      const selId = getResortsId(results[1]);

      expect(results.length).toBe(3);
      expect(selId).toBe('PRODUCT-BBB');
    });

    it('getSelected() should return the selected Entity', () => {
      const result = ResortsSelectors.getSelected(state) as ResortsEntity;
      const selId = getResortsId(result);

      expect(selId).toBe('PRODUCT-BBB');
    });

    it('getResortsLoaded() should return the current "loaded" status', () => {
      const result = ResortsSelectors.getResortsLoaded(state);

      expect(result).toBe(true);
    });

    it('getResortsError() should return the current "error" state', () => {
      const result = ResortsSelectors.getResortsError(state);

      expect(result).toBe(ERROR_MSG);
    });
  });
});
