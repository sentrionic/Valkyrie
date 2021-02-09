import { ApiProperty } from '@nestjs/swagger';
import { MemberResponse } from './MemberResponse';

export class DMChannelResponse {
  @ApiProperty({ type: String })
  id!: string;
  @ApiProperty({ type: MemberResponse })
  user!: MemberResponse;
}
