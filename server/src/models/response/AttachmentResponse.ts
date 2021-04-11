import { ApiProperty } from '@nestjs/swagger';

export class AttachmentResponse {
  @ApiProperty({ type: String })
  url!: string;
  @ApiProperty({ type: String })
  filetype!: string;
  @ApiProperty({ type: String })
  filename!: string;
}
