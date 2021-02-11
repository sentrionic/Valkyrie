import { ApiProperty } from '@nestjs/swagger';

export class ChannelInput {
  @ApiProperty({ description: 'Channel Name. 3 to 30 characters' })
  name!: string;

  @ApiProperty({ required: false, default: true })
  isPublic?: boolean = true;

  @ApiProperty({
    type: [String],
    required: false,
    description: 'Member IDs that are allowed in the channel'
  })
  members?: string[] = [];
}
