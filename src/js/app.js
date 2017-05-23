import React from 'react';
import ReactDOM from 'react-dom';
import { applyMiddleware, createStore, combineReducers } from 'redux';
import logger from 'redux-logger';
import { connect, Provider } from 'react-redux';

import App from './containers/App';
import roomActivityReducer, { add as addMessage, clear as clearMessages } from './stores/RoomActivityStore';
import roomsReducer, { set as setRooms } from './stores/RoomsStore';
import sessionReducer from './stores/SessionStore';
import Transport, { createReducer, listRooms, resolveWsPath } from './transports/JSONSocket';
import '../scss/app.scss';

// transport layer for server
const server = new Transport(resolveWsPath(window.location, '/api.v1'));

// prepare app to connect to react store
// TODO: only apply logger if NODE_ENV is "development"
const store = createStore(
  combineReducers({
    roomActivities: roomActivityReducer,
    rooms: roomsReducer,
    session: sessionReducer,
    transport: createReducer(server),
  }),
  undefined,
  applyMiddleware(logger),
);
const mapStateToProps = (state) => {
  const { roomActivities, rooms } = state;
  return { roomActivities, rooms };
};
const ConnectedApp = connect(mapStateToProps)(App);

// subscribe server broadcast
server.subscribe((message) => {
  const { entity = '', data = {} } = message;
  if (entity === 'roomActivities') {
    const { action } = data;
    switch (action) {
      case 'message': {
        const { message: messageText, user_id: userID } = data;
        store.dispatch(addMessage({
          message: messageText,
          userID,
        }));
        break;
      }
      default: {
        // do nothing
        console.log('TomatoRPG: received unknown roomActivities', data);
      }
    }
  } else if (message.type === 'response' && message.status === 'success') {
    if (message.entity === 'rooms') {
      switch (message.action) {
        case 'list': {
          // either create, update, list or delete rooms
          console.log(message.data);
          if (Array.isArray(message.data)) {
            store.dispatch(setRooms(message.data));
          }
          break;
        }
        default:
          // clear current room messages
          store.dispatch(clearMessages());
      }
    } else {
      console.log(`TomatoRPG: ${message.entity}.${message.action} ${message.status}`);
    }
  } else if (message.type === 'response' && message.status === 'error') {
    // TODO: throw error and somehow handles it
    console.error(`TomatoRPG: ${message.entity}.${message.action} ${message.status}: ${message.error}`);
  } else {
    console.log('TomatoRPG: received unknown server message', message);
  }
});

// initialize server connection transport
server.connect(() => {
  // only on first connection, not on re-connect
  store.dispatch(listRooms());
});

ReactDOM.render(
  <Provider store={store}>
    <ConnectedApp />
  </Provider>,
  document.getElementById('app'),
);
