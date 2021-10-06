// eslint-disable-next-line max-classes-per-file
export interface LoginDTO {
  email: string;
  password: string;

  [key: string]: any;
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
