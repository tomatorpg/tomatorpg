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
        <ol>
          { rooms.map((room, index) => {
            const key = `room-${index}`;
            return (
              <li key={key} className="room">
                {room.name} <button type="button" onClick={() => this.joinRoom(room.id)}>Join</button>
              </li>
            );
          }) }
        </ol>
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
