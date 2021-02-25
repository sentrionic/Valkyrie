import socketIOClient from 'socket.io-client';

export const getSocket = () => socketIOClient(`${process.env.REACT_APP_API!}/ws`, {
  transports: ['websocket'],
  upgrade: false
});
