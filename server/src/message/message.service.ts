import { BadRequestException, Inject, Injectable, NotFoundException, UnauthorizedException } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { Message } from '../entities/message.entity';
import { User } from '../entities/user.entity';
import { Channel } from '../entities/channel.entity';
import { uploadToS3 } from '../utils/fileUtils';
import { BufferFile } from '../types/BufferFile';
import { MessageResponse } from '../models/response/MessageResponse';
import { MessageInput } from '../models/dto/MessageInput';
import { SocketService } from '../socket/socket.service';

@Injectable()
export class MessageService {
  constructor(
    @InjectRepository(User) private userRepository: Repository<User>,
    @InjectRepository(Message) private messageRepository: Repository<Message>,
    @InjectRepository(Channel) private channelRepository: Repository<Channel>,
    private readonly socketService: SocketService
  ) {}

  async getMessages(
    channelId: string,
    userId: string,
    cursor?: string | null,
  ): Promise<MessageResponse[]> {
    const channel = await this.channelRepository.findOne({
      where: { id: channelId },
    });

    // if (!channel.public) {
    //   const member = await this.pcMemberRepository.findOne({
    //     where: { channelId, userId },
    //   });
    //
    //   if (!member) {
    //     throw new Error('Not Authorized');
    //   }
    // }

    const queryBuilder = this.messageRepository
      .createQueryBuilder('message')
      .leftJoinAndSelect('message.user', 'user')
      .where('message."channelId" = :channelId', { channelId })
      .orderBy('message.createdAt', 'DESC')
      .limit(35);

    if (cursor) {
      const timeString = new Date(cursor).getTime().toString();
      const time = timeString.substring(0, timeString.length - 3);
      queryBuilder.andWhere('message."createdAt" < (to_timestamp(:cursor))', {
        cursor: time,
      });
    }

    const messages = await queryBuilder.getMany();
    return messages.map(m => m.toJSON());
  }

  async createMessage(
    userId: string,
    channelId: string,
    input: MessageInput,
    file?: BufferFile,
  ): Promise<boolean> {

    if (!file && !input.text) {
      throw new BadRequestException();
    }

    const message = this.messageRepository.create({ ...input });

    if (file) {
      const directory = `channels/${channelId}`;
      const url = await uploadToS3(
        directory,
        file
      );
      message.filetype = 'image/webp';
      message.url = url;
    }

    const user = await this.userRepository.findOneOrFail({ where: { id: userId } });

    const channel = await this.channelRepository.findOneOrFail({
      where: { id: channelId },
    });

    message.user = user;
    message.channel = channel;

    await message.save();

    this.socketService.sendMessage({ room: channelId, message: message.toJSON() });

    return true;
  }

  async editMessage(
    userId: string,
    id: string,
    text: string,
  ): Promise<boolean> {

      let message = await this.messageRepository.findOneOrFail({
        where: { id },
        relations: ['user', 'channel'],
      });

      if (!message) {
        throw new NotFoundException();
      }

      if (message.user.id !== userId) {
        throw new UnauthorizedException();
      }

      await this.messageRepository.update(id, { text });

      message = await this.messageRepository.findOneOrFail({
        where: { id },
        relations: ['user', 'channel'],
      });

      this.socketService.editMessage({ room: message.channel.id, message: message.toJSON() });

      return true;
  }

  async deleteMessage(userId: string, id: string): Promise<boolean> {
    const message: Message = await this.messageRepository.findOneOrFail({
      where: { id },
      relations: ['user', 'channel'],
    });

    if (!message) {
      throw new NotFoundException();
    }

    if (message.user.id !== userId) {
      throw new UnauthorizedException();
    }

    const deleteId = message.id;

    await this.messageRepository.remove(message);

    message.id = deleteId;

    this.socketService.deleteMessage({ room: message.channel.id, message: message.toJSON() });

    return true;
  }

}
