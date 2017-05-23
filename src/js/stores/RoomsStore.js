
export function set(rooms) {
  return {
    type: 'ROOMS_SET',
    rooms,
  };
}

export default function reducer(state = [], action) {
  const { type } = action;
  switch (type) {
    case 'ROOMS_SET': {
      const { rooms } = action;
      return Array.from(rooms);
    }
    default:
      // do nothing
  }
  return state;
}
