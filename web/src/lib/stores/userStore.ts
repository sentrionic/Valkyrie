import create from 'zustand';
import { persist } from 'zustand/middleware';
import { Account } from '../models/account';

type AccountState = {
  current: Account | null;
  setUser: (account: Account) => void;
  logout: () => void;
};

export const userStore = create<AccountState>(
  persist(
    /* eslint-disable */
    (set, _) => ({
      current: null,
      setUser: (account) => set({ current: account }),
      logout: () => set({ current: null }),
    }),
    {
      name: 'user-storage',
    }
  )
);
