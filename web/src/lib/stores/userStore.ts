import create from 'zustand';
import { persist } from 'zustand/middleware';

import { AccountResponse } from '../api/models';

type AccountState = {
  current: AccountResponse | null;
  setUser: (account: AccountResponse) => void;
  logout: () => void;
};

export const userStore = create<AccountState>(
  persist(
    (set) => ({
      current: null,
      setUser: (account) => set({ current: account }),
      logout: () => set({ current: null }),
    }),
    {
      name: 'user-storage',
    }
  )
);
