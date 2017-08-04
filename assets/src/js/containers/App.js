/* eslint
react/forbid-prop-types: 'warn'
*/
import React from 'react';
import { connect } from 'react-redux';
import { BrowserRouter as Router, Route } from 'react-router-dom';

import Nav from './Nav';
import Rooms from './Rooms';
import Room from './Room';
import RoomCharacters from './RoomCharacters';
import Character from './Character';
import { setRoomID } from '../stores/SessionStore';
import { clear as clearMessages } from '../stores/RoomActivityStore';
import { createRoom, joinRoom, listRooms, listRoomActivities } from '../transports/JSONSocket';

const ConnectedNav = connect(
  (state) => {
    const { session: { user } } = state;
    return { user };
  },
)(Nav);

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
      dispatch(listRoomActivities(roomID));
    },
  }),
)(Room);

const ConnectedRoomCharacters = connect(
  (state) => {
    const { roomCharacters } = state;
    return { roomCharacters };
  },
  dispatch => ({
    dispatch,
  }),
)(RoomCharacters);

const ConnectedCharacter = connect(
  (state) => {
    const { session } = state;
    console.log('session', session);
    return { session };
  },
  dispatch => ({
    dispatch,
  }),
)(Character);

const App = () => (
  <Router>
    <div>
      <ConnectedNav />
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
          render={({ match, history }) =>
          (<ConnectedRoom
            roomID={match.params.roomID}
            listCharacters={() => history.push(`/rooms/${match.params.roomID}/characters`)}
          />)}
        />
        <Route
          exact
          path="/rooms/:roomID/characters"
          render={({ match, history }) =>
          (<ConnectedRoomCharacters
            roomID={match.params.roomID}
            onClose={() => history.push(`/rooms/${match.params.roomID}`)}
            createHandler={() => history.push(`/rooms/${match.params.roomID}/characters/create`)}
            closeHandler={() => history.push(`/rooms/${match.params.roomID}`)}
          />)}
        />
        <Route
          exact
          path="/rooms/:roomID/characters/create"
          render={({ match, history }) =>
          (<ConnectedCharacter
            roomID={match.params.roomID}
            postSubmit={() => history.push(`/rooms/${match.params.roomID}/characters`)}
            closeHandler={() => history.push(`/rooms/${match.params.roomID}/characters`)}
          />)}
        />
      </main>
    </div>
  </Router>
);

export default App;
