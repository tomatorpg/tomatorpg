import reducer, { add, set } from '../../src/js/stores/RoomCharactersStore';

test('reducer to implement ROOM_CHARACTERS_ADD', () => {
  expect(reducer(undefined, add({
    name: 'hello',
  }))).toEqual([
    { name: 'hello' },
  ]);

  expect(reducer([
    { name: 'world' },
  ], add({
    name: 'hello',
  }))).toEqual([
    { name: 'world' },
    { name: 'hello' },
  ]);
});

test('reducer to implement ROOM_CHARACTERS_SET', () => {
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
