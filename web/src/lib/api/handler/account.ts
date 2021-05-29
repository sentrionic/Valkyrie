import { AxiosResponse } from 'axios';
import { AccountResponse, Member, RequestResponse } from '../models';
import { request } from '../setupAxios';

export const getAccount = (): Promise<AxiosResponse<AccountResponse>> => request.get('/account');

export const updateAccount = (body: FormData): Promise<AxiosResponse<AccountResponse>> =>
  request.put('/account', body, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });

export const getFriends = (): Promise<AxiosResponse<Member[]>> => request.get('/account/me/friends');

export const getPendingRequests = (): Promise<AxiosResponse<RequestResponse[]>> => request.get('/account/me/pending');

export const sendFriendRequest = (id: string): Promise<AxiosResponse<boolean>> => request.post(`/account/${id}/friend`);

export const acceptFriendRequest = (id: string): Promise<AxiosResponse<boolean>> =>
  request.post(`/account/${id}/friend/accept`);

export const declineFriendRequest = (id: string): Promise<AxiosResponse<boolean>> =>
  request.post(`/account/${id}/friend/cancel`);

export const removeFriend = (id: string): Promise<AxiosResponse<boolean>> => request.delete(`/account/${id}/friend`);
