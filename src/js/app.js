import React from 'react';
import ReactDOM from 'react-dom';
import { applyMiddleware, createStore, combineReducers } from 'redux';
import logger from 'redux-logger';
import { connect, Provider } from 'react-redux';

import App from './containers/App';
import roomActivityReducer, { add as addMessage } from './stores/RoomActivityStore';
import roomsReducer from './stores/RoomsStore';
import sessionReducer from './stores/SessionStore';
import Transport, { createReducer, resolveWsPath } from './transports/JSONSocket';
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
  const { roomActivities } = state;
  return { roomActivities };
};
const ConnectedApp = connect(mapStateToProps)(App);

// subscribe server broadcast
server.subscribe((message) => {
  const { entity = '', data = {} } = message;
  if (entity === 'roomActivities') {
    const { action } = data;
    switch (action) {
      case 'message': {
        store.dispatch(addMessage(data.message));
        break;
      }
      default: {
        // do nothing
        console.log('TomatoRPG: received unknown roomActivities', data);
      }
    }
  } else if (message.type === 'response' && message.status === 'success') {
    console.log(`TomatoRPG: ${message.entity}.${message.action} ${message.status}`);
  } else if (message.type === 'response' && message.status === 'error') {
    // TODO: throw error and somehow handles it
    console.error(`TomatoRPG: ${message.entity}.${message.action} ${message.status}: ${message.error}`);
  } else {
    console.log('TomatoRPG: received unknown server message', message);
  }
});

// initialize server connection transport
server.connect();

ReactDOM.render(
  <Provider store={store}>
    <ConnectedApp />
  </Provider>,
  document.getElementById('app'),
);
