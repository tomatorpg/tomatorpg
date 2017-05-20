const defaultState = [];

export function add(message) {
  return {
    type: 'MESSAGE_ADD',
    message,
  };
}

const reducer = (state = defaultState, action = {}) => {
  switch (action.type) {
    case 'MESSAGE_ADD': {
      const { message } = action;
      const activities = Array.from(state);
      activities.push({
        type: 'message',
        message,
      });
      return activities;
    }

    default:
      // placeholder
  }

  return state;
};

export default reducer;
