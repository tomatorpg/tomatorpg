/* eslint
react/forbid-prop-types: 'warn'
*/

import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import PropTypes from 'prop-types';

class Rooms extends Component {

  createRoom() {
    const { createRoom } = this.props;
    createRoom();
  }

  joinRoom(id) {
    console.log(`goto /rooms/${id}`);
    this.context.router.history.push(`/rooms/${id}`);
  }

  render() {
    const { rooms = [] } = this.props;
    return (rooms.length > 0) ? (
      <div id="rooms">
        <button type="button" onClick={evt => this.createRoom(evt)}>Create</button>
        <ul className="rooms">
          { rooms.map((room) => {
            const key = `room-${room.id}`;
            const roomDisplayName = (typeof room.name === 'string' && room.name.trim() !== '') ?
              room.name : `Room ${room.id}`;
            return (
              <li key={key} className="room">
                <Link className="room-name" to={`/rooms/${room.id}`}>
                  {roomDisplayName}
                </Link>
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

Rooms.contextTypes = {
  router: PropTypes.object,
};

Rooms.propTypes = {
  createRoom: PropTypes.func,
  rooms: PropTypes.array,
};

Rooms.defaultProps = {
  createRoom: () => {},
  rooms: [],
};

export default Rooms;
