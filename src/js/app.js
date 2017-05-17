import React from 'react';
import ReactDOM from 'react-dom';
import Hello from './containers/Hello';
import '../scss/app.scss';

ReactDOM.render(
  <Hello />,
  document.getElementById('hello')
);

const chatbox = document.getElementById('chatbox');
const messages = document.getElementById('messages');

function localWebSocketPath(path) {
  var loc = window.location, new_uri;
  if (loc.protocol === "https:") {
      new_uri = "wss:";
  } else {
      new_uri = "ws:";
  }
  new_uri += "//" + loc.host;
  new_uri += path;
  return new_uri;
}

var exampleSocket = new WebSocket(localWebSocketPath('/api.v1'), "protocolOne");
exampleSocket.onopen = function (event) {
  exampleSocket.send(JSON.stringify({
    'action': 'sign_in',
  }));
};
exampleSocket.onmessage = function (event) {
  let broadcast = JSON.parse(event.data);
  let messageBox = document.createElement('div');
  messageBox.appendChild(document.createTextNode(broadcast.message));
  messages.appendChild(messageBox);
  messages.scrollTop += 50;
}
exampleSocket.onclose = function (event) {
  console.info('connection closed');
}

chatbox.addEventListener('submit', function (e) {
  e.preventDefault();
  exampleSocket.send(JSON.stringify({
    'message': e.target.children.to_send.value,
  }));
  console.log('form submitted: ', e.target.children.to_send.value);
  e.target.children.to_send.value = ''; // clear message
})
