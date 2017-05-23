import reducer, { add, clear } from '../../src/js/stores/RoomActivityStore';

test('reducer to implement ROOM_ACTIVITIES_MESSAGE', () => {
  expect(reducer(undefined, add('some message'))).toEqual(
    [
      {
        type: 'message',
        message: 'some message',
      },
    ],
  );
});

test('reducer to implement ROOM_ACTIVITIES_CLEAR', () => {
  expect(reducer([{ type: 'message', message: 'hello' }], clear())).toEqual([]);
});
