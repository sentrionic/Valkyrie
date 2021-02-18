import { ApiProperty } from '@nestjs/swagger';

export class ChannelResponse {
  @ApiProperty({ type: String })
  id!: string;

  @ApiProperty({ type: String })
  name!: string;

  @ApiProperty()
  isPublic!: boolean;

  @ApiProperty({ type: String })
  createdAt!: string;

  @ApiProperty({ type: String })
  updatedAt!: string;

  @ApiProperty()
  hasNotification?: boolean;
}
