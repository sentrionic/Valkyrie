import React, { useEffect } from 'react';
import { userStore } from '../../lib/stores/userStore';
import { getSocket } from '../../lib/api/getSocket';

export const GlobalState: React.FC = ({ children }) => {

  const current = userStore(state => state.current);

  useEffect(() => {
    if (current) {

      const disconnect = () => {
        socket.emit('toggleOffline');
        socket.disconnect();
      }

      const socket = getSocket();
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
