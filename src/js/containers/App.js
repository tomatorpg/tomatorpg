/* eslint
react/forbid-prop-types: 'warn'
*/
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { BrowserRouter as Router, Route } from 'react-router-dom';

import Nav from './Nav';
import Rooms from './Rooms';
import Room from './Room';
import { setRoomID } from '../stores/SessionStore';
import { clear as clearMessages } from '../stores/RoomActivityStore';
import { createRoom, joinRoom, listRooms, replayRoom } from '../transports/JSONSocket';

const ConnectedRooms = connect(
  (state) => {
    const { rooms } = state;
    return { rooms };
  },
  dispatch => ({
    dispatch,
    createRoom: () => {
      dispatch(createRoom());
      dispatch(listRooms());
    },
  }),
)(Rooms);

const ConnectedRoom = connect(
  (state) => {
    const { roomActivities } = state;
    return { roomActivities };
  },
  dispatch => ({
    dispatch,
    onLoad: (props) => {
      const { roomID } = props;
      // dispatch these events when loading the room
      dispatch(clearMessages());
      dispatch(joinRoom(roomID));
      dispatch(setRoomID(roomID)); // TODO:need to check if room exists
      dispatch(replayRoom(roomID));
    },
  }),
)(Room);

class App extends Component {
  render() {
    return (
      <Router>
        <div>
          <Nav />
          <main>
            <Route
              exact
              path="/"
              render={() => <ConnectedRooms />}
            />
            <Route
              exact
              path="/rooms"
              render={() => <ConnectedRooms />}
            />
            <Route
              path="/rooms/:roomID"
              render={({ match }) =>
              (<ConnectedRoom
                roomID={match.params.roomID}
              />)}
            />
          </main>
        </div>
      </Router>
    );
  }
}

export default App;
