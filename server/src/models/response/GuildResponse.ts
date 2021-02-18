import { ApiProperty } from '@nestjs/swagger';

export class GuildResponse {
  @ApiProperty({ type: String })
  id!: string;
  @ApiProperty({ type: String })
  name!: string;
  @ApiProperty({ type: String })
  icon?: string;
  @ApiProperty({ type: String })
  default_channel_id!: string;
  @ApiProperty({ type: String })
  ownerId!: string;
  @ApiProperty({ type: String })
  createdAt!: string;
  @ApiProperty({ type: String })
  updatedAt!: string;
  @ApiProperty({ type: Boolean })
  hasNotification!: boolean;
}
