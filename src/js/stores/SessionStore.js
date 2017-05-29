const defaultState = {
  roomID: '',
  token: '',
  user: {
    name: 'Visitor',
  },
};

export function setRoomID(roomID) {
  return {
    type: 'SESSION_ROOM_ID_SET',
    roomID,
  };
}

export function setToken(token) {
  return {
    type: 'SESSION_TOKEN_SET',
    token,
  };
}

export function setUser(user) {
  return {
    type: 'SESSION_USER_SET',
    user,
  };
}

const reducer = (state = defaultState, action = {}) => {
  switch (action.type) {
    case 'SESSION_ROOM_ID_SET': {
      const { roomID } = action;
      return Object.assign(
        {},
        state,
        {
          roomID,
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
    default:
      // placeholder
  }

  return state;
};

export default reducer;
