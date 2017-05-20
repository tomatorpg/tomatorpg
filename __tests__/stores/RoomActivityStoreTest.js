import reducer, { add } from '../../src/js/stores/RoomActivityStore';

test('reducer to implement MESSAGE_ADD', () => {
  expect(reducer(undefined, add('some message'))).toEqual(
    [
      {
        type: 'message',
        message: 'some message',
      },
    ],
  );
});
