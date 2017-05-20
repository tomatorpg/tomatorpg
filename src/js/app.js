import React from 'react';
import ReactDOM from 'react-dom';

import App from './containers/App';
import Transport, { resolveWsPath } from './transports/JSONSocket';
import '../scss/app.scss';

const server = new Transport(resolveWsPath(window.location, '/api.v1'));

server.subscribe((message) => {
  console.log(`server receive new message: ${message.message}`);

  // TODO: dispatch redux action to update messages store
});
server.connect();

ReactDOM.render(
  <App server={server} />,
  document.getElementById('app'),
);
