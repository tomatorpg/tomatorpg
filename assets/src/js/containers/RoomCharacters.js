/* eslint
react/forbid-prop-types: 'warn',
*/

import React, { Component } from 'react';
import PropTypes from 'prop-types';
import Modal from 'react-modal';

class RoomCharacters extends Component {
  render() {
    const { roomCharacters, closeHandler, createHandler } = this.props;
    return (
      <Modal
        isOpen
        contentLabel="character-form"
        overlayClassName="modal-overlay characters-list"
      >
        <div className="modal-actions">
          <button type="button" onClick={() => closeHandler()}>Close</button>
        </div>
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
      </Modal>
    );
  }
}

RoomCharacters.propTypes = {
  roomCharacters: PropTypes.array,
  closeHandler: PropTypes.func,
  createHandler: PropTypes.func,
};

RoomCharacters.defaultProps = {
  roomCharacters: [],
  closeHandler: () => {},
  createHandler: () => {},
};

export default RoomCharacters;
