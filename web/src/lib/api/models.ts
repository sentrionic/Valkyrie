interface BaseModel {
  id: string;
  createdAt: string;
  updatedAt: string;
}

export interface FieldError {
  field: string;
  message: string;
}

export interface Member extends BaseModel {
  username: string;
  image: string;
  isOnline: boolean;
  isFriend: boolean;
  nickname?: string | null;
  color?: string | null;
}

export interface Message extends BaseModel {
  text?: string;
  filetype?: string;
  url?: string;
  user: Member;
}

export interface AccountResponse extends BaseModel {
  username: string;
  email: string;
  image: string;
}

export interface Channel extends BaseModel {
  name: string;
  isPublic: boolean;
}

export interface Guild extends BaseModel {
  name: string;
  ownerId: string;
  default_channel_id: string;
}

export interface DMChannel {
  id: string;
  user: Member;
}

export interface RequestResponse {
  id: string;
  username: string;
  image: string;
  type: number;
}
