import { ApiProperty } from '@nestjs/swagger';

export class MemberResponse {
  @ApiProperty({ type: String })
  id!: string;
  @ApiProperty({ type: String })
  username!: string;
  @ApiProperty({ type: String })
  image!: string;
  @ApiProperty({ type: Boolean })
  isOnline!: boolean;
  @ApiProperty({ type: String })
  createdAt!: string;
  @ApiProperty({ type: String })
  updatedAt!: string;
  @ApiProperty({ type: Boolean })
  isFriend!: boolean;
  @ApiProperty({ type: String, required: false })
  nickname?: string;
  @ApiProperty({ type: String, required: false })
  color?: string;
}
