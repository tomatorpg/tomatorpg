import React from 'react';
import ReactDOM from 'react-dom';

import Hello from './containers/Hello';
import Transport, { resolveWsPath } from './transports/JSONSocket';
import '../scss/app.scss';

ReactDOM.render(
  <Hello />,
  document.getElementById('hello'),
);

const chatbox = document.getElementById('chatbox');
const messages = document.getElementById('messages');

const server = new Transport(resolveWsPath(window.location, '/api.v1'));

server.subscribe((message) => {
  const messageBox = document.createElement('div');
  messageBox.appendChild(document.createTextNode(message.message));
  messages.appendChild(messageBox);
  messages.scrollTop += 50;
});
server.connect();

chatbox.addEventListener('submit', (e) => {
  e.preventDefault();
  server.send({
    message: e.target.children.to_send.value,
  });
  console.log('form submitted: ', e.target.children.to_send.value);
  e.target.children.to_send.value = ''; // clear message
});
