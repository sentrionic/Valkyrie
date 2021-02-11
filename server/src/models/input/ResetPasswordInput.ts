import { ApiProperty } from '@nestjs/swagger';

export class ResetPasswordInput {
  @ApiProperty({
    type: String,
    description: 'The from the email provided token.',
  })
  token!: string;
  @ApiProperty({ type: String, description: 'Min 6, max 150 characters.' })
  newPassword!: string;
  @ApiProperty({
    type: String,
    description: 'Must be the same as the newPassword value.',
  })
  confirmNewPassword!: string;
}

export class ForgotPasswordInput {
  @ApiProperty({
    type: String,
    description: 'User Email.',
  })
  email!: string;
}