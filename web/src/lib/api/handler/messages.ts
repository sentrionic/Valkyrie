import { AxiosResponse } from 'axios';
import { Message } from '../models';
import { request } from '../setupAxios';

export const getMessages = (id: string, cursor?: string): Promise<AxiosResponse<Message[]>> =>
  request.get(`messages/${id}/${cursor ? `?cursor=${cursor}` : ''}`);

export const sendMessage = (
  channelId: string,
  data: FormData,
  onUploadProgress?: (e: any) => void
): Promise<AxiosResponse<void>> =>
  request.post(`messages/${channelId}`, data, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
    onUploadProgress,
  });

export const deleteMessage = (id: string): Promise<AxiosResponse<boolean>> => request.delete(`messages/${id}`);

export const editMessage = (id: string, text: string): Promise<AxiosResponse<boolean>> =>
  request.put(`messages/${id}`, { text });
