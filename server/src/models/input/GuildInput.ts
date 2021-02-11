import { ApiProperty } from '@nestjs/swagger';

export class GuildInput {
  @ApiProperty({ description: 'Guild Name. 3 to 30 characters' })
  name!: string;
}
