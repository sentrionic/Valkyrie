import create from 'zustand';

type FriendStoreType = {
  isPending: boolean;
  toggleDisplay: () => void;
};

export const friendStore = create<FriendStoreType>(
  (set, get) => ({
    isPending: false,
    toggleDisplay: () => set(state => ({ isPending: !state.isPending })),
  })
);
