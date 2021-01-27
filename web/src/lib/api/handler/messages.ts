import { AxiosResponse } from 'axios';
import { Message } from '../models';
import { request } from '../setupAxios';

export const getMessages = (id: string, cursor?: string): Promise<AxiosResponse<Message[]>> =>
  request.get(`channels/${id}/messages${cursor ? `?cursor=${cursor}` : ''}`);

export const sendMessage = (channelId: string, data: FormData): Promise<AxiosResponse<boolean>> =>
  request.post(`channels/${channelId}/messages`, data, {
    headers: {
      "Content-Type": "multipart/form-data",
    },
  });