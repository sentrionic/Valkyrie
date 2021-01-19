import Axios, { AxiosResponse } from "axios";
import {
  ChangePasswordInput,
  LoginDTO,
  RegisterDTO,
  ResetPasswordInput,
} from "./dtos/models";
import { AccountResponse } from "./response/accountresponse";

const request = Axios.create({
  baseURL: process.env.REACT_APP_API_URL,
  withCredentials: true,
});

export const register = (
  body: RegisterDTO
): Promise<AxiosResponse<AccountResponse>> =>
  request.post("/account/register", body);

export const login = (
  body: LoginDTO
): Promise<AxiosResponse<AccountResponse>> =>
  request.post("/account/login", body);

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
