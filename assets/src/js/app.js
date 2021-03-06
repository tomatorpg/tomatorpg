import React from 'react';
import ReactDOM from 'react-dom';
import { applyMiddleware, createStore, combineReducers } from 'redux';
import logger from 'redux-logger';
import { Provider } from 'react-redux';

import App from './containers/App';
import roomActivityReducer, { add as addMessage } from './stores/RoomActivityStore';
import roomCharactersReducer, { add as addCharacter } from './stores/RoomCharactersStore';
import roomsReducer, { set as setRooms } from './stores/RoomsStore';
import sessionReducer, { setUser } from './stores/SessionStore';
import Transport, { createReducer, whoami, joinRoom, listRooms, resolveWsPath } from './transports/JSONSocket';
import '../scss/app.scss';

// transport layer for server
const server = new Transport(resolveWsPath(window.location, '/api.v1'));

// prepare app to connect to react store
// TODO: only apply logger if NODE_ENV is "development"
const store = createStore(
  combineReducers({
    roomActivities: roomActivityReducer,
    roomCharacters: roomCharactersReducer,
    rooms: roomsReducer,
    session: sessionReducer,
    transport: createReducer(server),
  }),
  undefined,
  applyMiddleware(logger),
);

// join the previous room on re-connect
server.subscribe('open', () => {
  const state = store.getState();
  if (state.session.roomID !== '') {
    server.dispatch(joinRoom(state.session.roomID));
  }
});

// subscribe server broadcast
server.subscribe('message', (message) => {
  const {
    entity = '',
    method = '',
    message_type: messageType = '',
    status = '',
    data = {},
  } = message;

  if (messageType === 'broadcast') {
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
          console.log('TomatoRPG: received unknown roomActivities', message);
        }
      }
    }
  } else if (messageType === 'response' && status === 'success') {
    if (entity === 'rooms' && method === 'list') {
      if (Array.isArray(message.data)) {
        store.dispatch(setRooms(message.data));
      }
    } else if (entity === '' && method === 'whoami') {
      store.dispatch(setUser(message.data));
    } else if (entity === 'roomActivities' && method === 'list') {
      console.log('roomActivities.list = data', data);
      if (Array.isArray(data)) {
        for (let i = 0; i < data.length; i += 1) {
          const roomActivity = data[i];
          const { action } = roomActivity;
          switch (action) {
            case 'message': {
              const { message: messageText, user_id: userID } = roomActivity;
              store.dispatch(addMessage({
                message: messageText,
                userID,
              }));
              break;
            }
            case 'createCharacter': {
              const { meta: character } = roomActivity;
              store.dispatch(addCharacter(character));
              break;
            }
            default: {
              // do nothing
              console.log('TomatoRPG: received unknown roomActivities', message);
            }
          }
        }
      }
    } else if (entity === 'roomActivities' && method === 'create') {
      console.log('roomActivities.list = data', data);
    } else {
      console.log(`TomatoRPG: ${message.entity}.${message.method} ${message.status}`);
    }
  } else if (message.message_type === 'response' && message.status === 'error') {
    // TODO: throw error and somehow handles it
    console.error(`TomatoRPG: ${message.entity}.${message.method} ${message.status}: ${message.error}`);
  } else {
    console.log('TomatoRPG: received unknown server message', message);
  }
});

// initialize server connection transport
server.connect(() => {
  // only on first connection, not on re-connect
  store.dispatch(listRooms());
  store.dispatch(whoami());
});

ReactDOM.render(
  <Provider store={store}>
    <App />
  </Provider>,
  document.getElementById('app'),
);
