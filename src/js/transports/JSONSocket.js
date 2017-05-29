// type name to mimic redux store action
const JSONSocketCmd = 'TOMATO_RPC';

/**
 * ping
 * @param {string} message string to send in the room
 * @return {Object} action suitable to dispatch to pseudo reducer and server
 */
export function ping() {
  return {
    type: JSONSocketCmd,
    tomatorpc: '0.1',
    context: 'session',
    method: 'ping',
  };
}

/**
 * messageInRoom
 * @param {string} message string to send in the room
 * @return {Object} action suitable to dispatch to pseudo reducer and server
 */
export function messageInRoom(message) {
  // TODO: Support sending message as character
  return {
    type: JSONSocketCmd,
    tomatorpc: '0.1',
    context: 'session',
    entity: 'roomActivities',
    method: 'create',
    payload: {
      action: 'message',
      message,
    },
  };
}

/**
 * createCharInRoom
 * @param {string} name of the character to create
 * @return {Object} action suitable to dispatch to pseudo reducer and server
 */
export function createCharInRoom(name) {
  return {
    type: JSONSocketCmd,
    tomatorpc: '0.1',
    context: 'session',
    entity: 'roomActivities',
    method: 'create',
    payload: {
      action: 'createCharacter',
      name,
    },
  };
}

/**
 * joinRoom
 * @param {string} roomID for the room to join
 * @return {Object} action suitable to dispatch to pseudo reducer and server
 */
export function joinRoom(roomID) {
  return {
    type: JSONSocketCmd,
    tomatorpc: '0.1',
    context: 'session',
    entity: 'rooms',
    room_id: roomID,
    method: 'join',
  };
}

/**
 * replayRoom
 * @param {string} roomID for the room to join
 * @return {Object} action suitable to dispatch to pseudo reducer and server
 */
export function replayRoom(roomID) {
  return {
    type: JSONSocketCmd,
    tomatorpc: '0.1',
    context: 'session',
    entity: 'rooms',
    room_id: roomID,
    method: 'replay',
  };
}

/**
 * createRoom
 * @param {string} name of the room to create
 * @return {Object} action suitable to dispatch to pseudo reducer and server
 */
export function createRoom(name = '(no name)') {
  // TODO: Support sending message / chat as character
  return {
    type: JSONSocketCmd,
    tomatorpc: '0.1',
    entity: 'rooms',
    method: 'create',
    name,
  };
}

/**
 * listRoom
 * @param {string} name of the room to create
 * @return {Object} action suitable to dispatch to pseudo reducer and server
 */
export function listRooms() {
  // TODO: Support sending message / chat as character
  return {
    type: JSONSocketCmd,
    tomatorpc: '0.1',
    entity: 'rooms',
    method: 'list',
  };
}

/**
 * deleteRoom
 * @param {string} roomID for the room to join
 * @return {Object} action suitable to dispatch to pseudo reducer and server
 */
export function deleteRoom(id) {
  return {
    type: JSONSocketCmd,
    tomatorpc: '0.1',
    entity: 'rooms',
    id,
    method: 'delete',
  };
}

/**
 * createReducer
 * @param {string} roomID for the room to join
 * @return {function} redux compatible reducer function that don't
 *                    really chnage state.
 */
export function createReducer(socket) {
  return (state = true, action) => {
    const { type } = action;
    if (type === JSONSocketCmd) {
      socket.dispatch(action);
    }
    return true;
  };
}

/**
 * resolveWsPathhelper
 * @param {window.location | url.URL} current window location or any valid url
 * @param {string} wsPath websocket entity path
 * @return {string} uri string of the websocket path
 */
export function resolveWsPath(uri, wsPath) {
  const protocol = (uri.protocol === 'https:') ? 'wss:' : 'ws:';
  return `${protocol}//${uri.host}${wsPath}`;
}

/**
 * @class JSONSocket
 * @description transport layer implementation to generic WebSocket RPC
 */
class JSONSocket {

  constructor(uri) {
    this.uri = uri;
    this.subscribers = {
      message: [],
      open: [],
    };
    this.ready = false;
    this.pending = [];
  }

  dispatch(action) {
    if (this.ready) {
      this.webSocket.send(JSON.stringify(action));
    } else {
      console.log('socket not ready. push the action into pending state.');
      this.pending.push(action);
    }
  }

  subscribe(eventName, subscriber) {
    if (typeof subscriber !== 'function') {
      throw new Exception('subscriber is not a function');
    }
    if (typeof this.subscribers[eventName] === 'undefined') {
      throw new Exception(`JSONSocket has no ${eventName} event`);
    }
    this.subscribers[eventName].push(subscriber);
  }

  connect(callback) {
    console.info(`%cTomatoRPG transport%c: %cconnecting %c${this.uri}`, 'font-weight: bold', 'color: inherit', 'color: rgb(200,100,0)', 'color: red');
    this.webSocket = new WebSocket(this.uri, 'tomatorpc-v1');
    this.webSocket.onopen = () => {
      console.info('%cTomatoRPG transport%c: %cconnected.', 'font-weight: bold', 'color: inherit', 'color: green');
      this.ready = true;
      this.dispatch(ping());
      if (typeof callback === 'function') {
        callback();
      }
      while (this.pending.length) {
        const action = this.pending.shift();
        this.dispatch(action);
      }
      for (let i = 0; i < this.subscribers.open.length; i += 1) {
        (async () => this.subscribers.open[i]())();
      }
    };
    this.webSocket.onmessage = (evt) => {
      try {
        const broadcast = JSON.parse(evt.data);
        for (let i = 0; i < this.subscribers.message.length; i += 1) {
          (async () => this.subscribers.message[i](broadcast))();
        }
      } catch (err) {
        console.groupCollapsed('Unable to parse broadcast message');
        console.error('error', err);
        console.error('raw broadcast event:', evt);
        console.groupEnd();
      }
    };
    this.webSocket.onclose = () => {
      console.info('%cTomatoRPG transport%c: %cconnection closed. reconnect...', 'font-weight: bold', 'color: inherit', 'color: red');
      this.ready = false;
      window.setTimeout(() => {
        this.connect();
      }, 1000);
    };
  }

}

export default JSONSocket;
