import create from 'zustand';
import { persist } from 'zustand/middleware';

type SettingsState = {
  showMembers: boolean;
  toggleShowMembers: () => void;
};

export const settingsStore = create(
  persist<SettingsState>(
    (set, get) => ({
      showMembers: true,
      toggleShowMembers: () => set({ showMembers: !get().showMembers }),
    }),
    {
      name: 'settings-storage',
    }
  )
);
