import { AxiosResponse } from 'axios';
import { request } from '../setupAxios';
import { DMChannel } from '../models';

export const getUserDMs = (): Promise<AxiosResponse<DMChannel[]>> =>
  request.get("/channels/me/dm");
