import { AxiosResponse } from 'axios';
import { Channel, Guild, Member } from '../models';
import { request } from '../setupAxios';
import { GuildDto } from '../dtos/GuildDto';
import { InviteDto } from '../dtos/InviteDto';
import { ChannelInput } from '../dtos/ChannelInput';

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

export const getGuildMembers = (id: string): Promise<AxiosResponse<Member[]>> =>
  request.get(`guilds/${id}/members`);

export const createChannel = (id: string, input: ChannelInput): Promise<AxiosResponse<boolean>> =>
  request.post(`channels/${id}`, input);

export const leaveGuild = (id: string): Promise<AxiosResponse<boolean>> =>
  request.delete(`guilds/${id}`);

export const editChannel = (guildId: string, channelId: string, input: ChannelInput): Promise<AxiosResponse<boolean>> =>
  request.put(`channels/${guildId}/${channelId}`, input);


export const deleteChannel = (guildId: string, channelId: string): Promise<AxiosResponse<boolean>> =>
  request.delete(`channels/${guildId}/${channelId}`);

export const getPrivateChannelMembers = (channelId: string): Promise<AxiosResponse<string[]>> =>
  request.get(`channels/${channelId}/members`);
