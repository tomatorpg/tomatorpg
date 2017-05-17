import reducer from '../../src/js/stores/ActivityStore';

test('reducer to implement SAY', () => {
  expected(reducer(undefined, {
    type: 'SAY',
    message: 'some message',
  })).toEqual({
    activities: [
      {
        type: 'SAY',
        message: 'some message',
      },
    ],
  });
});
