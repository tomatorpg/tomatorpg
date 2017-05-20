import reducer, { setUser, setToken } from '../../src/js/stores/SessionStore';

test('reducer to implement SESSION_USER_SET', () => {
  expect(reducer(undefined, setUser({
    name: 'Someone',
  })).user).toEqual({
    name: 'Someone',
  });
});

test('reducer to implement SESSION_TOKEN_SET', () => {
  expect(reducer(undefined, setToken('some-token')).token)
    .toEqual('some-token');
});
