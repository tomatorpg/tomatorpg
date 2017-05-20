const defaultState = {
  user: {
    name: 'Visitor',
  },
  token: '',
};

export function setUser(user) {
  return {
    type: 'SESSION_USER_SET',
    user,
  };
}

export function setToken(token) {
  return {
    type: 'SESSION_TOKEN_SET',
    token,
  };
}

const reducer = (state = defaultState, action = {}) => {
  switch (action.type) {
    case 'SESSION_USER_SET': {
      const { user: { name } } = action;
      return Object.assign(
        {},
        state,
        {
          user: {
            name,
          },
        },
      );
    }
    case 'SESSION_TOKEN_SET': {
      const { token } = action;
      return Object.assign(
        {},
        state,
        {
          token,
        },
      );
    }
    default:
      // placeholder
  }

  return state;
};

export default reducer;
