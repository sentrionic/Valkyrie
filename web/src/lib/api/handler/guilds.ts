import { AxiosResponse } from 'axios';
import { Guild, Member } from '../models';
import { request } from '../setupAxios';
import { GuildInput } from '../dtos/GuildInput';
import { InviteInput } from '../dtos/InviteInput';
import { GuildMemberInput } from "../dtos/GuildMemberInput";

export const getUserGuilds = (): Promise<AxiosResponse<Guild[]>> =>
  request.get("/guilds");

export const createGuild = (input: GuildInput): Promise<AxiosResponse<Guild>> =>
  request.post("guilds/create", input);

export const joinGuild = (input: InviteInput): Promise<AxiosResponse<Guild>> =>
  request.post("guilds/join", input);

export const getInviteLink = (id: string): Promise<AxiosResponse<string>> =>
  request.get(`guilds/${id}/invite`);

export const getGuildMembers = (id: string): Promise<AxiosResponse<Member[]>> =>
  request.get(`guilds/${id}/members`);

export const leaveGuild = (id: string): Promise<AxiosResponse<boolean>> =>
  request.delete(`guilds/${id}`);

export const editGuild = (id: string, input: FormData): Promise<AxiosResponse<boolean>> =>
  request.put(`guilds/${id}`, input, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });

export const deleteGuild = (id: string): Promise<AxiosResponse<boolean>> =>
  request.delete(`guilds/${id}/delete`);

export const getGuildMemberSettings = (id: string): Promise<AxiosResponse<GuildMemberInput>> =>
  request.get(`guilds/${id}/member`);

export const changeGuildMemberSettings = (id: string, input: GuildMemberInput): Promise<AxiosResponse<boolean>> =>
  request.put(`guilds/${id}/member`, input);

export const getBanList = (id: string): Promise<AxiosResponse<Member[]>> =>
  request.get(`guilds/${id}/bans`);

export const kickMember = (guildId: string, memberId: string): Promise<AxiosResponse<boolean>> =>
  request.post(`guilds/${guildId}/kick`, { memberId });

export const banMember = (guildId: string, memberId: string): Promise<AxiosResponse<boolean>> =>
  request.post(`guilds/${guildId}/ban`, { memberId });

export const unbanMember = (guildId: string, memberId: string): Promise<AxiosResponse<boolean>> =>
  request.delete(`guilds/${guildId}/ban`, { data: { memberId } });
