import reducer, { set } from '../../src/js/stores/RoomsStore';

test('reducer to implement ROOMS_SET', () => {
  expect(reducer(undefined, set([
    { name: 'hello' },
  ]))).toEqual([
    { name: 'hello' },
  ]);

  expect(reducer([
    { name: 'world' },
  ], set([
    { name: 'hello' },
  ]))).toEqual([
    { name: 'hello' },
  ]);
});
