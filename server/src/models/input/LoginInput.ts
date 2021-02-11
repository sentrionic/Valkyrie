import { User } from '../../entities/user.entity';
import { ApiProperty } from '@nestjs/swagger';

export class LoginInput implements Partial<User> {
  @ApiProperty({
    type: String,
    description: 'Must be a valid email.'
  })
  email!: string;

  @ApiProperty({ type: String, description: 'Min 6, max 150 characters.' })
  password!: string;
}
