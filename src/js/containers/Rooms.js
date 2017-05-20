/* eslint
react/forbid-prop-types: 'warn'
*/

import React, { Component } from 'react';
import PropTypes from 'prop-types';

class Rooms extends Component {
  render() {
    const { rooms = [] } = this.props;
    return (rooms.length > 0) ? (
      <div id="rooms">
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
  rooms: PropTypes.array,
};

Rooms.defaultProps = {
  rooms: [],
};

export default Rooms;
