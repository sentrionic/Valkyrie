import { Injectable, NotFoundException, UnauthorizedException } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { Message } from '../entities/message.entity';
import { User } from '../entities/user.entity';
import { Channel } from '../entities/channel.entity';
import { uploadToS3 } from '../utils/fileUtils';
import { BufferFile } from '../types/BufferFile';
import { MessageData } from '../types/MessageData';

@Injectable()
export class MessageService {
  constructor(
    @InjectRepository(Message) private messageRepository: Repository<Message>,
    @InjectRepository(Channel) private channelRepository: Repository<Channel>,
  ) {}

  async getMessages(
    cursor: string,
    channelId: string,
    userId: string,
  ): Promise<Message[]> {
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

    return await queryBuilder.getMany();
  }

  async createMessage(
    user: User,
    channelId: string,
    text?: string,
    file?: BufferFile,
  ): Promise<boolean> {
    try {
      const messageData: MessageData = { text };

      if (file) {
        const directory = `channels/${channelId}`;
        const url = await uploadToS3(
          directory,
          file
        );
        messageData.filetype = 'image/webp';
        messageData.url = url;
      }

      const channel = await this.channelRepository.findOne({
        where: { id: channelId },
      });

      const message = this.messageRepository.create();

      //TODO: Improve
      message.text = text;
      message.url = messageData.url;
      message.filetype = messageData.filetype;
      message.user = user;
      message.channel = channel;

      await message.save();

      // const asyncFunc = async () => {
      //   this.pubsub.publish(MESSAGE_SUBSCRIPTION, {
      //     channelId,
      //     channelMessage: {
      //       message: {
      //         message,
      //         operation: MessageOperation.NEW,
      //       },
      //     },
      //   });
      // };
      //
      // asyncFunc();

      return true;
    } catch (err) {
      console.log(err);
      return false;
    }
  }

  async editMessage(
    user: User,
    id: string,
    text: string,
  ): Promise<boolean> {

      let message: Message = await this.messageRepository.findOne({
        where: { id },
        relations: ['user', 'channel'],
      });

      if (!message) {
        throw new NotFoundException();
      }

      if (message.user.id !== user.id) {
        throw new UnauthorizedException();
      }

      message = await this.messageRepository.save({ ...message, text });

      // const asyncFunc = async () => {
      //   this.pubsub.publish(MESSAGE_SUBSCRIPTION, {
      //     channelId: message.channel.id,
      //     channelMessage: {
      //       message: {
      //         message,
      //         operation: MessageOperation.EDIT,
      //       },
      //     },
      //   });
      // };
      //
      // asyncFunc();

      return true;
  }

  async deleteMessage(userId: string, id: string): Promise<boolean> {
    const message: Message = await this.messageRepository.findOne({
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

    // const asyncFunc = async () => {
    //   this.pubsub.publish(MESSAGE_SUBSCRIPTION, {
    //     channelId: message.channel.id,
    //     channelMessage: {
    //       message: {
    //         message,
    //         operation: MessageOperation.DELETE,
    //       },
    //     },
    //   });
    // };
    //
    // asyncFunc();

    return true;
  }

}
