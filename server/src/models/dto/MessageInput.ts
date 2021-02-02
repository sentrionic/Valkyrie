import { ApiProperty } from '@nestjs/swagger';

export class MessageInput {
  @ApiProperty({
    type: String,
    required: false,
    description: "The message. Must not be empty"
  })
  text?: string;

  @ApiProperty({ type: String, required: false })
  file?: string;
}
