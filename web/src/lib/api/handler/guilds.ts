import { AxiosResponse } from 'axios';
import { Guild, Member } from '../models';
import { request } from '../setupAxios';
import { GuildInput } from '../dtos/GuildInput';
import { InviteInput } from '../dtos/InviteInput';

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

export const editGuild = (id: string, input: GuildInput): Promise<AxiosResponse<boolean>> =>
  request.put(`guilds/${id}`, input);

export const deleteGuild = (id: string): Promise<AxiosResponse<boolean>> =>
  request.delete(`guilds/${id}/delete`);
