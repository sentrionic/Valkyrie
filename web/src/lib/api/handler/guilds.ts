import { AxiosResponse } from 'axios';
import { Channel, Guild } from '../models';
import { request } from '../setupAxios';
import { GuildDto } from '../dtos/GuildDto';
import { InviteDto } from '../dtos/InviteDto';

export const getUserGuilds = (): Promise<AxiosResponse<Guild[]>> =>
  request.get("/guilds");

export const createGuild = (input: GuildDto): Promise<AxiosResponse<Guild>> =>
  request.post("guilds/create", input);

export const joinGuild = (input: InviteDto): Promise<AxiosResponse<Guild>> =>
  request.post("guilds/join", input);

export const getChannels = (id: string): Promise<AxiosResponse<Channel[]>> =>
  request.get(`channels/${id}`);

export const getInviteLink = (id: string): Promise<AxiosResponse<string>> =>
  request.get(`guilds/${id}/invite`);