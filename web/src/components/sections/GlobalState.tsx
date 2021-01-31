import React, { useEffect } from 'react';
import { userStore } from '../../lib/stores/userStore';
import socketIOClient from 'socket.io-client';

export const GlobalState: React.FC = ({ children }) => {

  const current = userStore(state => state.current);

  useEffect(() => {
    if (current) {

      const disconnect = () => {
        socket.emit('toggleOffline');
        socket.disconnect();
      }

      const socket = socketIOClient(process.env.REACT_APP_API_WS!);
      socket.emit('toggleOnline');

      window.addEventListener('beforeunload', disconnect);

      return () => disconnect();
    }
  }, [current]);

  return (
    <>
      {children}
    </>
  );
};