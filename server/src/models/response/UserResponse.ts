import { ApiProperty } from '@nestjs/swagger';

export class UserResponse {
  @ApiProperty({ type: String })
  id!: string;
  @ApiProperty({ type: String })
  username!: string;
  @ApiProperty({ type: String })
  email!: string;
  @ApiProperty({ type: String })
  image!: string;
  @ApiProperty({ type: String })
  createdAt!: string;
  @ApiProperty({ type: String })
  updatedAt!: string;
}
