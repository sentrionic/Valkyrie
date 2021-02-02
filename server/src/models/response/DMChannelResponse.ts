import { ApiProperty } from '@nestjs/swagger';

export class DMChannelResponse {
  @ApiProperty({ type: String })
  id!: string;
  @ApiProperty({ type: String })
  name!: string;
}
