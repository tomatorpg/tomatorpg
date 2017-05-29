import reducer, { setRoomID, setToken, setUser } from '../../src/js/stores/SessionStore';

test('reducer to implement SESSION_ROOM_ID_SET', () => {
  expect(reducer(undefined, setRoomID('some-id')).roomID)
    .toEqual('some-id');
});

test('reducer to implement SESSION_TOKEN_SET', () => {
  expect(reducer(undefined, setToken('some-token')).token)
    .toEqual('some-token');
});

test('reducer to implement SESSION_USER_SET', () => {
  expect(reducer(undefined, setUser({
    name: 'Someone',
  })).user).toEqual({
    name: 'Someone',
  });
});
