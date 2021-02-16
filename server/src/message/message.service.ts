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
  ) {
  }

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
    let time: string;
    if (cursor) {
      const timeString = new Date(cursor).getTime().toString();
      time = timeString.substring(0, timeString.length - 3);
    }

    const manager = getManager();
    const results = await manager.query(
      `
          SELECT "message".id,
                  "message".text,
                  "message".filetype,
                  "message".url,
                  "message"."createdAt",
                  "message"."updatedAt",
                  "user"."id" as "userId",
                  "user"."createdAt" as "ucreatedAt",
                  "user"."updatedAt" as "uupdatedAt",
                  "user"."username",
                  "user"."image",
                  "user"."isOnline",
                  ${!channel.dm ? "member.nickname, member.color," : ''}
                  exists(
                          select 1
                          from users
                                   left join friends f on users.id = f."user"
                          where f."friend" = "message"."userId"
                            and f."user" = $2
                      ) as "isFriend"
          FROM "messages" "message"
                   LEFT JOIN "users" "user" ON "user"."id" = "message"."userId"
                   ${!channel.dm ? 'LEFT JOIN members member on "message"."userId" = member."userId"' : ''}
          WHERE message."channelId" = $1 
          ${!channel.dm ? `AND member."guildId" = ${channel.guild.id}::text` : ''}
          ${cursor ? `AND message."createdAt" < (to_timestamp(${time}))` : ``}
          ORDER BY "message"."createdAt" DESC
              LIMIT 35
      `,
      [channelId, userId],
    );

    const messages: MessageResponse[] = [];
    results.map(m => messages.push({
      id: m.id,
      text: m.text,
      filetype: m.filetype,
      url: m.url,
      createdAt: m.createdAt,
      updatedAt: m.updatedAt,
      user: {
        id: m.userId,
        username: m.username,
        image: m.image,
        isOnline: m.isOnline,
        createdAt: m.ucreatedAt,
        updatedAt: m.uupdatedAt,
        isFriend: m.isFriend,
        nickname: m.nickname,
        color: m.color
      }
    }));

    return messages;
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

    const member = await this.memberRepository.findOne({
      where: {
        userId,
        guildId: channel.guild.id
      }
    });

    const response = message.toJSON(userId);
    response.user.nickname = member?.nickname;
    response.user.color = member?.color;

    this.socketService.sendMessage({ room: channelId, message: response });

    if (channel.dm) {
      // Open the DM and push it to the top
      getManager().query(
        `
            update dm_members
            set "isOpen" = true,
                "updatedAt" = CURRENT_TIMESTAMP
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

  /**
   * Checks if the current user is a member of that channel or dm.
   * Throws an UnauthorizedException if that's not the case
   * @param channel
   * @param userId
   * @private
   */
  private async isChannelMember(channel: Channel, userId: string): Promise<void> {
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
  }
}
