import { HttpException, HttpStatus, Injectable, NotFoundException, UnauthorizedException } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { User } from '../entities/user.entity';
import { Repository } from 'typeorm';
import { RegisterInput } from '../models/dto/RegisterInput';
import * as md5 from 'md5';
import e from 'express';
import { LoginInput } from '../models/dto/LoginInput';
import { v4 } from 'uuid';
import { redis } from '../config/redis';
import { FORGET_PASSWORD_PREFIX } from '../utils/constants';
import { sendEmail } from '../utils/sendEmail';
import { ResetPasswordInput } from '../models/dto/ResetPasswordInput';
import { ChangePasswordInput } from '../models/dto/ChangePasswordInput';
import { UpdateInput } from '../models/dto/UpdateInput';
import { BufferFile } from '../types/BufferFile';
import { uploadAvatarToS3 } from '../utils/fileUtils';
import { UserResponse } from '../models/response/UserResponse';
import * as argon2 from 'argon2';

@Injectable()
export class UserService {
  constructor(
    @InjectRepository(User) private userRepository: Repository<User>,
  ) {
  }

  async register(
    credentials: RegisterInput,
    req: e.Request,
  ): Promise<UserResponse> {
    const { email } = credentials;

    const emailTaken = await this.userRepository.findOne({
      where: { email },
    });

    if (emailTaken) {
      throw new HttpException(
        {
          errors: [
            {
              field: 'email',
              message: 'Email must be unique.',
            },
          ],
        },
        HttpStatus.BAD_REQUEST,
      );
    }

    const user = this.userRepository.create({ ...credentials, password: await argon2.hash(credentials.password) });
    user.image = `https://gravatar.com/avatar/${md5(email)}?d=identicon`;
    await user.save();

    req!.session!['userId'] = user.id;

    return user.toJSON();
  }

  async login(
    credentials: LoginInput,
    req: e.Request,
  ): Promise<UserResponse> {

    const { email, password } = credentials;

    const user = await this.userRepository.findOne({
      where: { email },
    });

    if (!user) {
      throw new NotFoundException();
    }

    const valid = await argon2.verify(user.password, password);

    if (!valid) {
      throw new UnauthorizedException();
    }

    req!.session!['userId'] = user.id;

    return user.toJSON();
  }

  async forgotPassword(email: string): Promise<boolean> {
    const user = await this.userRepository.findOne({ email });
    if (!user) {
      // the email is not in the db
      return true;
    }

    const token = v4();

    await redis.set(
      FORGET_PASSWORD_PREFIX + token,
      user.id,
      'ex',
      1000 * 60 * 60 * 24 * 3,
    ); // 3 days
    await sendEmail(
      email,
      `<a href='${process.env.CORS_ORIGIN}/reset-password/${token}'>Reset Password</a>`,
    );

    return true;
  }

  async resetPassword(
    input: ResetPasswordInput,
    req: e.Request,
  ): Promise<UserResponse> {
    const { newPassword, token } = input;

    const key = FORGET_PASSWORD_PREFIX + token;
    const userId = await redis.get(key);
    if (!userId) {
      throw new HttpException(
        {
          errors: [
            {
              field: 'token',
              message: 'Token expired',
            },
          ],
        },
        HttpStatus.BAD_REQUEST,
      );
    }

    const user = await this.userRepository.findOne({
      where: { id: userId },
    });

    if (!user) throw new NotFoundException();

    user.password = await argon2.hash(newPassword);

    await user.save();

    await redis.del(key);

    // log in user after change password
    req!.session!['userId'] = user.id;

    return user.toJSON();
  }

  async changePassword(
    input: ChangePasswordInput,
    userId: string,
  ): Promise<boolean> {
    const { newPassword, currentPassword } = input;

    const user = await this.userRepository.findOne({ where: { id: userId } });

    if (!user) {
      throw new NotFoundException();
    }

    const valid = await argon2.verify(user.password, currentPassword);

    if (!valid) {
      throw new UnauthorizedException();
    }

    user.password = await argon2.hash(newPassword);

    await user.save();

    return true;
  }

  async findCurrentUser(id: string): Promise<UserResponse> {
    const user = await this.userRepository.findOne(id);
    if (!user) {
      throw new NotFoundException({
        message: 'An account with that username or email does not exist.',
      });
    }
    return user.toJSON();
  }

  async updateUser(
    id: string,
    data: UpdateInput,
    image?: BufferFile,
  ): Promise<UserResponse> {
    const { email } = data;

    const user = await this.userRepository.findOneOrFail(id);

    if (user.email !== email) {
      const checkUsername = await this.userRepository.findOne({ email });
      if (checkUsername) {
        throw new HttpException(
          {
            errors: [
              {
                field: 'email',
                message: 'Email must be unique.',
              },
            ],
          },
          HttpStatus.BAD_REQUEST,
        );
      }
    }

    if (image) {
      const directory = `valkyrie/users/${id}`;
      data.image = await uploadAvatarToS3(directory, image);
    }

    if (!data.image || data.image === '') data.image = user.image;

    await this.userRepository.update({ id: user.id }, data);

    return this.findCurrentUser(user.id);
  }
}
