import { MemberResponse } from './MemberResponse';
import { ApiProperty } from '@nestjs/swagger';

export class MessageResponse {
  @ApiProperty({ type: String })
  id!: string;
  @ApiProperty({ type: String })
  @ApiProperty({ type: String })
  text?: string | null;
  @ApiProperty({ type: String })
  url?: string | null;
  @ApiProperty({ type: String })
  filetype?: string | null;
  @ApiProperty({ type: MessageResponse })
  user!: MemberResponse;
  @ApiProperty({ type: String })
  createdAt!: string;
  @ApiProperty({ type: String })
  updatedAt!: string;
}
