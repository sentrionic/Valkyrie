import { AxiosResponse } from 'axios';
import { request } from '../setupAxios';
import { Account } from '../../models/account';
import { Member } from '../../models/member';
import { FriendRequest } from '../../models/friend';

export const getAccount = (): Promise<AxiosResponse<Account>> => request.get('/account');

export const updateAccount = (body: FormData): Promise<AxiosResponse<Account>> =>
  request.put('/account', body, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });

export const getFriends = (): Promise<AxiosResponse<Member[]>> => request.get('/account/me/friends');

export const getPendingRequests = (): Promise<AxiosResponse<FriendRequest[]>> => request.get('/account/me/pending');

export const sendFriendRequest = (id: string): Promise<AxiosResponse<boolean>> => request.post(`/account/${id}/friend`);

export const acceptFriendRequest = (id: string): Promise<AxiosResponse<boolean>> =>
  request.post(`/account/${id}/friend/accept`);

export const declineFriendRequest = (id: string): Promise<AxiosResponse<boolean>> =>
  request.post(`/account/${id}/friend/cancel`);

export const removeFriend = (id: string): Promise<AxiosResponse<boolean>> => request.delete(`/account/${id}/friend`);
