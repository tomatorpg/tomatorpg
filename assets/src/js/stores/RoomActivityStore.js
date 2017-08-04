const defaultState = [];

export function add(message) {
  return {
    type: 'ROOM_ACTIVITIES_MESSAGE',
    message,
  };
}

export function clear() {
  return {
    type: 'ROOM_ACTIVITIES_CLEAR',
  };
}

const reducer = (state = defaultState, action = {}) => {
  switch (action.type) {
    case 'ROOM_ACTIVITIES_MESSAGE': {
      const { message } = action;
      const activities = Array.from(state);
      activities.push({
        type: 'message',
        message,
      });
      return activities;
    }

    case 'ROOM_ACTIVITIES_CLEAR': {
      return [];
    }

    default:
      // placeholder
  }

  return state;
};

export default reducer;
