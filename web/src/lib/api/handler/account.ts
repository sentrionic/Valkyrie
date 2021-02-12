import { AxiosResponse } from 'axios';
import { AccountResponse, Member } from '../models';
import { request } from '../setupAxios';

export const getAccount = (): Promise<AxiosResponse<AccountResponse>> =>
  request.get("/account");

export const updateAccount = (
  body: FormData
): Promise<AxiosResponse<AccountResponse>> =>
  request.put("/account", body, {
    headers: {
      "Content-Type": "multipart/form-data",
    },
  });

export const getFriends = (): Promise<AxiosResponse<Member[]>> =>
  request.get("/account/me/friends");
