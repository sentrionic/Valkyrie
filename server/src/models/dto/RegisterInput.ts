import { User } from '../../entities/user.entity';

export class RegisterInput implements Partial<User> {
  username: string;
  email: string;
  password: string;
}