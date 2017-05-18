const defaultState = {
  activities: [],
};

const reducer = (state = defaultState, action = {}) => {
  switch (action.type) {
    case 'SAY': {
      const { type, message } = action;
      const activities = Array.from(state.activities);
      activities.push({
        type,
        message,
      });
      return Object.assign(
        {},
        state,
        { activities },
      );
    }

    default:
      // placeholder
  }

  return state;
};

export default reducer;
