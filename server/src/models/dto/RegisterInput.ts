import { User } from '../../entities/user.entity';
import { ApiProperty } from '@nestjs/swagger';

export class RegisterInput implements Partial<User> {
  @ApiProperty({
    type: String,
    description: 'Min 3, max 30 characters.',
  })
  username!: string;

  @ApiProperty({
    type: String,
    description: 'Unique. Must be a valid email.'
  })
  email!: string;

  @ApiProperty({ type: String, description: 'Min 6, max 150 characters.' })
  password!: string;
}
