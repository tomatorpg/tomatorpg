
export function add(character) {
  return {
    type: 'ROOM_CHARACTERS_ADD',
    character,
  };
}

export function set(characters) {
  return {
    type: 'ROOM_CHARACTERS_SET',
    characters,
  };
}

export default function reducer(state = [], action = {}) {
  switch (action.type) {
    case 'ROOM_CHARACTERS_ADD': {
      const { character } = action;
      const next = Array.from(state);
      next.push(character);
      return next;
    }

    case 'ROOM_CHARACTERS_SET': {
      const { characters } = action;
      return Array.from(characters);
    }

    default:
      // placeholder
  }
  return state;
}
