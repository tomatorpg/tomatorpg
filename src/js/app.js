import React from 'react';
import ReactDOM from 'react-dom';
import { applyMiddleware, createStore, combineReducers } from 'redux';
import logger from 'redux-logger';
import { connect, Provider } from 'react-redux';

import App from './containers/App';
import roomActivityReducer, { add as addMessage } from './stores/RoomActivityStore';
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
server.subscribe((broadcast) => {
  const { entity = '', data = {} } = broadcast;
  if (entity === 'roomActivities') {
    const { action } = data;
    switch (action) {
      case 'message': {
        const { message } = data;
        store.dispatch(addMessage(message));
        break;
      }
      default: {
        // do nothing
        console.log('TomatoRPG: received unknown roomActivities', data);
      }
    }
  } else {
    console.log('TomatoRPG: received unknown server broadcast', broadcast);
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
