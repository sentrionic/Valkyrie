import { BadRequestException, Injectable, NotFoundException, UnauthorizedException } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { getManager, Repository } from 'typeorm';
import { Message } from '../entities/message.entity';
import { User } from '../entities/user.entity';
import { Channel } from '../entities/channel.entity';
import { deleteFile, uploadToS3 } from '../utils/fileUtils';
import { BufferFile } from '../types/BufferFile';
import { MessageResponse } from '../models/response/MessageResponse';
import { MessageInput } from '../models/input/MessageInput';
import { SocketService } from '../socket/socket.service';
import { PCMember } from '../entities/pcmember.entity';
import { Member } from '../entities/member.entity';
import { DMMember } from '../entities/dmmember.entity';

@Injectable()
export class MessageService {
  constructor(
    @InjectRepository(User) private userRepository: Repository<User>,
    @InjectRepository(Message) private messageRepository: Repository<Message>,
    @InjectRepository(Channel) private channelRepository: Repository<Channel>,
    @InjectRepository(Member) private memberRepository: Repository<Member>,
    @InjectRepository(PCMember)
    private pcMemberRepository: Repository<PCMember>,
    @InjectRepository(DMMember)
    private dmMemberRepository: Repository<DMMember>,
    private readonly socketService: SocketService
  ) {}

  /**
   * Returns 35 messages for the given channel.
   * Requires channel access.
   * Uses the "createdAt" attribute as the cursor
   * @param channelId
   * @param userId
   * @param cursor
   */
  async getMessages(
    channelId: string,
    userId: string,
    cursor?: string | null,
  ): Promise<MessageResponse[]> {
    const channel = await this.channelRepository.findOne({
      where: { id: channelId },
      relations: ['guild']
    });

    await this.isChannelMember(channel, userId);

    const queryBuilder = this.messageRepository
      .createQueryBuilder('message')
      .leftJoinAndSelect('message.user', 'user')
      .leftJoinAndSelect('user.friends', 'friends')
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
    return messages.map(m => m.toJSON(userId));
  }

  async createMessage(
    userId: string,
    channelId: string,
    input: MessageInput,
    file?: BufferFile,
  ): Promise<boolean> {

    const channel = await this.channelRepository.findOneOrFail({
      where: { id: channelId },
      relations: ['guild']
    });

    await this.isChannelMember(channel, userId);

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
      message.filetype = file.mimetype;
      message.url = url;
    }

    message.user = await this.userRepository.findOneOrFail({ where: { id: userId }, relations: ['friends'] });
    message.channel = channel;

    await message.save();

    this.socketService.sendMessage({ room: channelId, message: message.toJSON(userId) });

    if (channel.dm) {
      // Open the DM and push it to the top
      getManager().query(
        `
        update dm_members set "isOpen" = true, "updatedAt" = CURRENT_TIMESTAMP 
        where "channelId" = $1 
        `, [channelId]
      );
      this.socketService.pushDMToTop({ room: channelId, channelId })
    }

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
        relations: ['user', 'channel', 'friends'],
      });

      this.socketService.editMessage({ room: message.channel.id, message: message.toJSON(userId) });

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

    if (message.url) {
      await deleteFile(message.url);
    }

    await this.messageRepository.remove(message);

    message.id = deleteId;

    this.socketService.deleteMessage({ room: message.channel.id, message: message.toJSON(userId) });

    return true;
  }

  private async isChannelMember(channel: Channel, userId: string): Promise<boolean> {
    // Check if user has access to private channel
    if (!channel.isPublic) {
      // Channel is DM -> Check if one of the members
      if (channel.dm) {
        const member = await this.dmMemberRepository.findOne({
          where: { channelId: channel.id, userId },
        });

        if (!member) {
          throw new UnauthorizedException('Not Authorized');
        }
        // Channel is private
      } else {
        const member = await this.pcMemberRepository.findOne({
          where: { channelId: channel.id, userId },
        });

        if (!member) {
          throw new UnauthorizedException('Not Authorized');
        }
      }
      // Check if user has access to the channel
    } else {
      const member = await this.memberRepository.findOneOrFail({
        where: { guildId: channel.guild.id, userId },
      });

      if (!member) {
        throw new UnauthorizedException('Not Authorized');
      }
    }

    return true;
  }

}
