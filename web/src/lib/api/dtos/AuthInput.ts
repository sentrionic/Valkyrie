export interface LoginDTO {
  [key: string]: any;
  email: string;
  password: string;
}

export interface RegisterDTO extends LoginDTO {
  username: string;
}

export class ChangePasswordInput {
  currentPassword!: string;
  newPassword!: string;
  confirmNewPassword!: string;
}

export class ResetPasswordInput {
  token!: string;
  newPassword!: string;
  confirmNewPassword!: string;
}
