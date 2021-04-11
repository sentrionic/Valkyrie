import { ApiProperty } from '@nestjs/swagger';
import { MemberResponse } from './MemberResponse';
import { AttachmentResponse } from './AttachmentResponse';

export class MessageResponse {
  @ApiProperty({ type: String })
  id!: string;
  @ApiProperty({ type: String })
  text?: string | null;
  @ApiProperty({ type: MessageResponse })
  user!: MemberResponse;
  @ApiProperty({ type: AttachmentResponse })
  attachment?: AttachmentResponse | null;
  @ApiProperty({ type: String })
  createdAt!: string;
  @ApiProperty({ type: String })
  updatedAt!: string;
}
