import { ApiProperty } from '@nestjs/swagger';

export class UpdateInput {
  @ApiProperty({ type: String, description: 'Unique. Must be a valid email.' })
  email!: string;

  @ApiProperty({
    type: String,
    description: 'Min 3, max 30 characters.',
  })
  username!: string;

  @ApiProperty({ type: String, required: false })
  image?: string;
}
