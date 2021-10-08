import { AxiosResponse } from 'axios';
import { request } from '../setupAxios';
import { ChannelInput } from '../dtos/ChannelInput';
import { Channel } from '../../models/channel';

export const getChannels = (id: string): Promise<AxiosResponse<Channel[]>> => request.get(`channels/${id}`);

export const createChannel = (id: string, input: ChannelInput): Promise<AxiosResponse<Channel>> =>
  request.post(`channels/${id}`, input);

export const editChannel = (channelId: string, input: ChannelInput): Promise<AxiosResponse<boolean>> =>
  request.put(`channels/${channelId}`, input);

export const deleteChannel = (channelId: string): Promise<AxiosResponse<boolean>> =>
  request.delete(`channels/${channelId}`);

export const getPrivateChannelMembers = (channelId: string): Promise<AxiosResponse<string[]>> =>
  request.get(`channels/${channelId}/members`);
