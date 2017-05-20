import React from 'react';
import ReactDOM from 'react-dom';
import { applyMiddleware, createStore, combineReducers } from 'redux';
import logger from 'redux-logger';
import { connect, Provider } from 'react-redux';

import App from './containers/App';
import RoomActivityReducer, { add as addMessage } from './stores/RoomActivityStore';
import Transport, { resolveWsPath } from './transports/JSONSocket';
import '../scss/app.scss';

// prepare app to connect to react store
// TODO: only apply logger if NODE_ENV is "development"
const store = createStore(
  combineReducers({
    roomActivities: RoomActivityReducer,
  }),
  undefined,
  applyMiddleware(logger),
);
const mapStateToProps = (state) => {
  const { roomActivities } = state;
  return { roomActivities };
};
const ConnectedApp = connect(mapStateToProps)(App);

// transport layer for server
const server = new Transport(resolveWsPath(window.location, '/api.v1'));
server.subscribe((message) => {
  // TODO: add message should address the character / user
  // or give the store another action for chat
  store.dispatch(addMessage(message.message));
});
server.connect();

ReactDOM.render(
  <Provider store={store}>
    <ConnectedApp server={server} />
  </Provider>,
  document.getElementById('app'),
);
