import { MemberResponse } from './MemberResponse';

export class MessageResponse {
  id: string;
  text?: string;
  url?: string;
  filetype: string;
  user: MemberResponse;
  createdAt: string;
  updatedAt: string;
}
