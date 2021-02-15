import { ApiProperty } from '@nestjs/swagger';
import { MemberResponse } from './MemberResponse';

export class MessageResponse {
  @ApiProperty({ type: String })
  id!: string;
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
