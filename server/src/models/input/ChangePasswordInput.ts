import { ApiProperty } from '@nestjs/swagger';

export class ChangePasswordInput {
  @ApiProperty({ type: String })
  currentPassword!: string;

  @ApiProperty({ type: String, description: 'Min 6, max 150 characters.' })
  newPassword!: string;

  @ApiProperty({
    type: String,
    description: 'Must be the same as the newPassword value.',
  })
  confirmNewPassword!: string;
}
