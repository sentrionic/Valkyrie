import { ApiProperty } from '@nestjs/swagger';

export class GuildMemberInput {
  @ApiProperty({ type: String, required: false, description: 'Min 3, max 30 characters.' })
  nickname?: string;

  @ApiProperty({ type: String, required: false, description: 'Must be a valid hex string' })
  color?: string;
}