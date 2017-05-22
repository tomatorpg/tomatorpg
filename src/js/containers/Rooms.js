/* eslint
react/forbid-prop-types: 'warn'
*/

import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { createRoom } from '../transports/JSONSocket';

class Rooms extends Component {

  createRoom() {
    const { dispatch } = this.props;
    dispatch(createRoom());
  }

  render() {
    const { rooms = [] } = this.props;
    return (rooms.length > 0) ? (
      <div id="rooms">
        <button type="button" onClick={evt => this.createRoom(evt)}>Create</button>
        <ul>
          { rooms.map((room, index) => {
            const key = `room-${index}`;
            return <li key={key} className="room">{room.name}</li>;
          }) }
        </ul>
      </div>
    ) : (
      <div id="rooms">
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
