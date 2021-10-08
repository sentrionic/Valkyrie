import { AxiosResponse } from 'axios';
import { GuildMemberInput } from '../dtos/GuildMemberInput';
import { request } from '../setupAxios';
import { Member } from '../../models/member';

export const getGuildMemberSettings = (id: string): Promise<AxiosResponse<GuildMemberInput>> =>
  request.get(`guilds/${id}/member`);

export const changeGuildMemberSettings = (id: string, input: GuildMemberInput): Promise<AxiosResponse<boolean>> =>
  request.put(`guilds/${id}/member`, input);

export const getBanList = (id: string): Promise<AxiosResponse<Member[]>> => request.get(`guilds/${id}/bans`);

export const kickMember = (guildId: string, memberId: string): Promise<AxiosResponse<boolean>> =>
  request.post(`guilds/${guildId}/kick`, { memberId });

export const banMember = (guildId: string, memberId: string): Promise<AxiosResponse<boolean>> =>
  request.post(`guilds/${guildId}/bans`, { memberId });

export const unbanMember = (guildId: string, memberId: string): Promise<AxiosResponse<boolean>> =>
  request.delete(`guilds/${guildId}/bans`, { data: { memberId } });
