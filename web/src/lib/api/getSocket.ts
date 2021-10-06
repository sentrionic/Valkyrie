import ReconnectingWebSocket from 'reconnecting-websocket';

export const getSocket = (): ReconnectingWebSocket => new ReconnectingWebSocket(process.env.REACT_APP_WS!);

let socket: ReconnectingWebSocket | null = null;
export const getSameSocket = (): ReconnectingWebSocket => {
  if (!socket) {
    socket = new ReconnectingWebSocket(process.env.REACT_APP_WS!);
  }

  return socket;
};
