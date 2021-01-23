import { AxiosResponse } from 'axios';
import { AccountResponse } from '../models';
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
