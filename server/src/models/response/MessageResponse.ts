import { MemberResponse } from './MemberResponse';

class MessageResponse {
  id: string;
  text?: string;
  url?: string;
  filetype: string;
  user: MemberResponse;
  createdAt: string;
  updatedAt: string;
}
