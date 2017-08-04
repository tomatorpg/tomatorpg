/* eslint
react/forbid-prop-types: 'warn',
*/

import React, { Component } from 'react';
import PropTypes from 'prop-types';
import Modal from 'react-modal';
import { createCharInRoom } from '../transports/JSONSocket';
import { defaultState as defaultSession } from '../stores/SessionStore';

class Character extends Component {
  submitHandler(evt) {
    const { dispatch, postSubmit, session } = this.props;
    evt.preventDefault();
    dispatch(createCharInRoom({
      name: this.nameInput.value,
      roomID: session.roomID,
      desc: this.descInput.value,
    }));
    postSubmit(); // run after dispatch
  }
  render() {
    return (
      <Modal
        isOpen
        contentLabel="character-form"
        overlayClassName="modal-overlay character-form"
      >
        <h2>Create Character</h2>
        <form onSubmit={evt => this.submitHandler(evt)}>
          <div className="field">
            <input
              id="character-name"
              type="textfield"
              placeholder="Name"
              ref={(input) => { this.nameInput = input; }}
            />
          </div>
          <div className="field">
            <textarea
              id="character-desc"
              placeholder="Description"
              ref={(input) => { this.descInput = input; }}
            />
          </div>
        </form>
      </Modal>
    );
  }
}

Character.propTypes = {
  dispatch: PropTypes.func,
  postSubmit: PropTypes.func,
  session: PropTypes.object,
};

Character.defaultProps = {
  dispatch: () => {},
  postSubmit: () => {},
  session: defaultSession,
};


export default Character;
