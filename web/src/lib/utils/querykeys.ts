export const gKey = 'guilds';
export const dmKey = 'dms';
export const aKey = 'account';
export const fKey = 'friends';
export const rKey = 'requests';
export const nKey = 'notification';

export const cKey = (guildId: string): string => `channels-${guildId}`;
export const mKey = (guildId: string): string => `members-${guildId}`;
export const vcKey = (guildId: string): string => `voice-members-${guildId}`;
