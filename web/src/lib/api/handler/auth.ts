import { AxiosResponse } from 'axios';
import { AccountResponse } from '../models';
import { request } from '../setupAxios';
import { ChangePasswordInput, LoginDTO, RegisterDTO, ResetPasswordInput } from '../dtos/AuthInput';

export const register = (
  body: RegisterDTO
): Promise<AxiosResponse<AccountResponse>> =>
  request.post("/account/register", body);

export const login = (
  body: LoginDTO
): Promise<AxiosResponse<AccountResponse>> =>
  request.post("/account/login", body);

export const logout = (): Promise<AxiosResponse> =>
  request.post("/account/logout");

export const forgotPassword = (
  email: string
): Promise<AxiosResponse<boolean>> =>
  request.post("/account/forgot-password", { email });

export const changePassword = (
  body: ChangePasswordInput
): Promise<AxiosResponse> => request.put("/account/change-password", body);

export const resetPassword = (
  body: ResetPasswordInput
): Promise<AxiosResponse<AccountResponse>> =>
  request.post("/account/reset-password", body);
