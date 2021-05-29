import create from 'zustand';

type HomeStoreType = {
  notifCount: number;
  requestCount: number;
  increment: () => void;
  setRequests: (r: number) => void;
  reset: () => void;
  resetRequest: () => void;
  isPending: boolean;
  toggleDisplay: () => void;
};

export const homeStore = create<HomeStoreType>((set, get) => ({
  notifCount: 0,
  requestCount: 0,
  increment: () => set((state) => ({ notifCount: state.notifCount + 1 })),
  reset: () => set((state) => ({ notifCount: 0 })),
  resetRequest: () => set((state) => ({ requestCount: 0 })),
  setRequests: (r) => set((state) => ({ requestCount: r })),
  isPending: false,
  toggleDisplay: () => set((state) => ({ isPending: !state.isPending })),
}));
