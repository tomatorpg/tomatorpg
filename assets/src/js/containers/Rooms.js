/* eslint
react/forbid-prop-types: 'warn'
*/

import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import PropTypes from 'prop-types';
import Notice from './Notice';

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
        <div className="container">
          <Notice>
            <div style={{ textAlign: 'center', marginBottom: '1em' }}>
              <img alt="logo" src="/assets/images/diceicon.png" style={{ width: '8%' }} />
            </div>
            <h1 style={{ textAlign: 'center' }}> Welcome to Tomato RPG. </h1>
            <p style={{ textAlign: 'center' }}>
              Enjoy your TRPG Game here.
            </p>
          </Notice>
          <div className="create-button">
            <button className="btn btn-primary" type="button" onClick={evt => this.createRoom(evt)}>Create</button>
          </div>
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
                    <button className="btn btn-primary" type="button" onClick={() => this.joinRoom(room.id)}>Join</button>
                  </div>
                </li>
              );
            }) }
          </ul>
        </div>
      </div>
    ) : (
      <div id="rooms">
        <div className="container">
          <Notice>
            <h1 style={{ textAlign: 'center' }}> Welcome to Tomato RPG. </h1>
            <p style={{ marginLeft: '1em', marginRight: '1em' }}> Lorem ipsum dolor sit amet, consectetur adipiscing elit.
              Duis nec dapibus nulla. Etiam eleifend risus leo,
              eu scelerisque lorem posuere vel.</p>
          </Notice>
          <button className="btn btn-primary" type="button" onClick={evt => this.createRoom(evt)}>Create</button>
          <p className="msg-no-room">There is no room yet.</p>
        </div>
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
