import { Member } from './member';

export interface Message {
  id: string;
  text?: string;
  createdAt: string;
  updatedAt: string;
  attachment?: Attachment;
  user: Member;
}

export interface Attachment {
  filename: string;
  filetype: string;
  url: string;
}
