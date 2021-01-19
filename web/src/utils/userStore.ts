import create from "zustand";
import { persist } from "zustand/middleware";
import { AccountResponse } from "../api/response/accountresponse";
import { getAccount } from "../api/setupAxios";

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
      name: "user-storage",
      getStorage: () => sessionStorage,
    }
  )
);

export const getCurrent = () => userStore((state) => state.current);
