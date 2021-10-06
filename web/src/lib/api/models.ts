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
  attachment?: Attachment;
  user: Member;
}

export interface Attachment extends BaseModel {
  filename: string;
  filetype: string;
  url: string;
}

export interface AccountResponse extends BaseModel {
  username: string;
  email: string;
  image: string;
}

export interface Channel extends BaseModel {
  name: string;
  isPublic: boolean;
  hasNotification?: boolean;
}

export interface Guild extends BaseModel {
  name: string;
  ownerId: string;
  default_channel_id: string;
  icon?: string;
  hasNotification?: boolean;
}

export interface DMChannel {
  id: string;
  user: Member;
}

export interface DMNotification extends DMChannel {
  count: number;
}

export interface RequestResponse {
  id: string;
  username: string;
  image: string;
  type: number;
}

export interface RouterProps {
  guildId: string;
  channelId: string;
}
