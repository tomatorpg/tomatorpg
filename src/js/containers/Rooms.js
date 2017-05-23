/* eslint
react/forbid-prop-types: 'warn'
*/

import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { createRoom, joinRoom } from '../transports/JSONSocket';

class Rooms extends Component {

  createRoom() {
    const { dispatch } = this.props;
    dispatch(createRoom());
  }

  joinRoom(id) {
    const { dispatch } = this.props;
    dispatch(joinRoom(id));
  }

  render() {
    const { rooms = [] } = this.props;
    return (rooms.length > 0) ? (
      <div id="rooms">
        <button type="button" onClick={evt => this.createRoom(evt)}>Create</button>
        <ul className="rooms">
          { rooms.map((room, index) => {
            const key = `room-${room.id}`;
            const roomDisplayName = (typeof room.name === 'string' && room.name.trim() !== '') ?
              room.name : `Room ${room.id}`;
            return (
              <li key={key} className="room">
                <div className="room-name">{roomDisplayName}</div>
                <div className="room-actions">
                  <button type="button" onClick={() => this.joinRoom(room.id)}>Join</button>
                </div>
              </li>
            );
          }) }
        </ul>
      </div>
    ) : (
      <div id="rooms">
        <button type="button" onClick={evt => this.createRoom(evt)}>Create</button>
        <p className="msg-no-room">There is no room yet.</p>
      </div>
    );
  }
}

Rooms.propTypes = {
  dispatch: PropTypes.func,
  rooms: PropTypes.array,
};

Rooms.defaultProps = {
  dispatch: () => {},
  rooms: [],
};

export default Rooms;
