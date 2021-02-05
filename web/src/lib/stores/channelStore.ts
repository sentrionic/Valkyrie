import create from 'zustand';

type ChannelState = {
  typing: string[];
  addTyping: (username: string) => void;
  removeTyping: (username: string) => void;
  reset: () => void;
};

export const channelStore = create<ChannelState>(
  (set, get) => ({
    typing: [],
    addTyping: (username) => set(state => ({ typing: [...state.typing, username] })),
    removeTyping: (username) => set(state => ({ typing: [...state.typing.filter(u => u !== username)] })),
    reset: () => set(state => ({ typing: [] }))
  })
);