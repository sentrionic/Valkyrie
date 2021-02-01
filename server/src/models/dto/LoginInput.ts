import { User } from '../../entities/user.entity';

export class LoginInput implements Partial<User> {
  email!: string;
  password!: string;
}