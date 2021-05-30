import ReconnectingWebSocket from 'reconnecting-websocket';

export const getSocket = () => new ReconnectingWebSocket(process.env.REACT_APP_WS!);

let socket: ReconnectingWebSocket | null = null;
export const getSameSocket = () => {
  if (!socket) {
    socket = new ReconnectingWebSocket(process.env.REACT_APP_WS!);
  }

  return socket;
};
