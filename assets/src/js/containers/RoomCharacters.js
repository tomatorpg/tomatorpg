/* eslint
react/forbid-prop-types: 'warn',
*/

import React, { Component } from 'react';
import PropTypes from 'prop-types';

class RoomCharacters extends Component {
  render() {
    const { roomCharacters, createHandler } = this.props;
    return (
      <div className="character-list">
        <h2>Characters</h2>
        <div>
          <button type="button" onClick={() => createHandler()}>Create</button>
        </div>
        <div>
          {roomCharacters.map(character => (
            <div className="character">
              <div className="name">{character.name}</div>
            </div>
          ))}
        </div>
      </div>
    );
  }
}

RoomCharacters.propTypes = {
  roomCharacters: PropTypes.array,
  createHandler: PropTypes.func,
};

RoomCharacters.defaultProps = {
  roomCharacters: [],
  createHandler: () => {},
};

export default RoomCharacters;
