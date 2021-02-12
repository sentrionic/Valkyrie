export const gKey = 'guilds';
export const dmKey = 'dms';
export const aKey = 'account';
export const fKey = 'friends';

export const cKey = (guildId: string) => `channels-${guildId}`;
export const mKey = (guildId: string) => `members-${guildId}`;
