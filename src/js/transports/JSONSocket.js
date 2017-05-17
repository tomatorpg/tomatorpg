export function resolveWsPath(uri, wsPath) {
  const protocol = (uri.protocol === 'https:') ? 'wss:' : 'ws:';
  return `${protocol}//${uri.host}${wsPath}`;
}

class JSONSocket {

  constructor(uri) {
    this.uri = uri;
    this.webSocket = new WebSocket(uri, 'protocolOne');
    this.subscribers = [];
  }

  send(payload) {
    this.webSocket.send(JSON.stringify(payload));
  }

  subscribe(subscriber) {
    if (typeof subscriber !== 'function') {
      throw new Exception('subscriber is not a function');
    }
    this.subscribers.push(subscriber);
  }

  connect() {
    console.log(`connecting ${this.uri}`);
    this.webSocket = new WebSocket(this.uri, 'protocolOne');
    this.webSocket.onopen = () => {
      this.webSocket.send(JSON.stringify({
        action: 'sign_in',
      }));
    };
    this.webSocket.onmessage = (event) => {
      const message = JSON.parse(event.data);
      for (let i = 0; i < this.subscribers.length; i += 1) {
        (async () => this.subscribers[i](message))();
      }
    };
    this.webSocket.onclose = () => {
      console.info('connection closed. reconnect.');
      window.setTimeout(() => {
        this.connect();
      }, 1000);
    };
  }

}

export default JSONSocket;
