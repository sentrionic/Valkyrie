import create from 'zustand';
import { persist } from 'zustand/middleware';

import { AccountResponse } from '../api/models';

type AccountState = {
  current: AccountResponse | null;
  setUser: (account: AccountResponse) => void;
  logout: () => void;
  isAuth: () => boolean;
};

export const userStore = create<AccountState>(
  persist(
    (set, get) => ({
      current: null,
      setUser: (account: AccountResponse) => set({ current: account }),
      logout: () => set({ current: null }),
      isAuth: () => get().current != null,
    }),
    {
      name: 'user-storage',
    }
  )
);
