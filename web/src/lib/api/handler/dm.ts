import { AxiosResponse } from 'axios';
import { request } from '../setupAxios';
import { DMChannel } from '../models';

export const getUserDMs = (): Promise<AxiosResponse<DMChannel[]>> => request.get('/channels/me/dm');

export const getOrCreateDirectMessage = (id: string): Promise<AxiosResponse<DMChannel>> =>
  request.post(`/channels/${id}/dm`);

export const closeDirectMessage = (id: string): Promise<AxiosResponse<boolean>> => request.delete(`/channels/${id}/dm`);
