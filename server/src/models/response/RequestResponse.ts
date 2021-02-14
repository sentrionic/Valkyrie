import { ApiProperty } from '@nestjs/swagger';

export class RequestResponse {
  @ApiProperty({ type: String })
  id!: string;
  @ApiProperty({ type: String })
  username!: string;
  @ApiProperty({ type: String })
  image!: string;
  @ApiProperty({ type: Number, description: "1: Incoming, 0: Outgoing" })
  type!: number;
}
