import { BadRequestException, Injectable, InternalServerErrorException, NotFoundException } from '@nestjs/common';
import { getManager, Repository } from 'typeorm';
import { Guild } from '../entities/guild.entity';
import { InjectRepository } from '@nestjs/typeorm';
import { Channel } from '../entities/channel.entity';
import { User } from '../entities/user.entity';
import { Member } from '../entities/member.entity';
import { nanoid } from 'nanoid';
import { redis } from '../config/redis';
import { INVITE_LINK_PREFIX } from '../utils/constants';
import { MemberResponse } from '../models/response/MemberResponse';
import { GuildResponse } from '../models/response/GuildResponse';
import { SocketService } from '../socket/socket.service';
import { idGenerator } from '../utils/idGenerator';

@Injectable()
export class GuildService {
  constructor(
    @InjectRepository(Guild) private guildRepository: Repository<Guild>,
    @InjectRepository(Channel) private channelRepository: Repository<Channel>,
    @InjectRepository(User) private userRepository: Repository<User>,
    @InjectRepository(Member) private memberRepository: Repository<Member>,
    private socketService: SocketService
  ) {
  }

  async getGuildMembers(guildId: string): Promise<MemberResponse[]> {
    const manager = getManager();
    return await manager.query(
      `select distinct u.id, u.username, u.image, u."isOnline", u."createdAt", u."updatedAt"
       from users as u
                join members m on u."id"::text = m."userId"
       where m."guildId" = $1
       order by u.username
      `,
      [guildId]
    );
  }

  async getUserGuilds(userId: string): Promise<GuildResponse[]> {
    const manager = getManager();
    return await manager.query(
      `select distinct g."id",
                       g."name",
                       g."ownerId",
                       g."createdAt",
                       g."updatedAt",
                       (select c.id as "default_channel_id"
                        from channels c
                                 join guilds g on g.id = c."guildId"
                        where g.id = member."guildId"
                        order by c."createdAt" limit 1)
       from guilds g
           join members as member
       on g."id"::text = member."guildId"
       where member."userId" = $1
       order by g."createdAt";`,
      [userId]
    );
  }

  async getDirectMessageMembers(
    guildId: string,
    userId: string
  ): Promise<User[]> {
    const manager = getManager();
    return await manager.query(
      `select distinct
       on (u.id) u.id, u.username
       from users as u
           join direct_message as dm
       on (u."id"::text = dm."senderId"::text)
           or (u."id"::text = dm."receiverId"::text)
       where ($1 = dm."senderId"::text
          or $1 = dm."receiverId"::text)
         and dm."guildId"::text = $2
         and u."id"::text != $1
      `,
      [userId, guildId]
    );
  }

  async createGuild(name: string, userId: string): Promise<GuildResponse> {
    try {
      let guild: Guild | null = null;
      let channel: Channel | null = null;

      await getManager().transaction(async (entityManager) => {
        guild = this.guildRepository.create({ ownerId: userId });
        channel = this.channelRepository.create({ name: 'general' });

        guild.name = name.trim();
        await guild.save();
        await entityManager.save(guild);

        channel.guild = guild;
        await channel.save();
        await entityManager.save(channel);

        await entityManager.insert(Member, {
          id: await idGenerator(),
          guildId: guild.id,
          userId,
        });
      });

      return {
        id: guild!.id,
        name: guild!.name,
        default_channel_id: channel!.id,
        createdAt: guild!.createdAt.toString(),
        updatedAt: guild!.updatedAt.toString()
      };
    } catch (err) {
      throw new InternalServerErrorException(err);
    }
  }

  async generateInviteLink(guildId: string): Promise<string> {
    const token = nanoid(8);
    await redis.set(INVITE_LINK_PREFIX + token, guildId, 'ex', 60 * 60 * 24); // 1 day expiration
    return `${process.env.CORS_ORIGIN}/${token}`;
  }

  async joinGuild(token: string, userId: string): Promise<GuildResponse> {
    if (token.includes('/')) {
      token = token.substring(token.lastIndexOf('/') + 1);
    }

    const guildId = await redis.get(INVITE_LINK_PREFIX + token);

    if (!guildId) {
      throw new NotFoundException();
    }

    const guild = await this.guildRepository.findOne(guildId);

    if (!guild) {
      throw new NotFoundException();
    }

    const isMember = await this.memberRepository.findOne({ where: { userId, guildId } });

    if (isMember) {
      throw new BadRequestException();
    }

    await this.memberRepository.insert({ id: await idGenerator(), userId, guildId });

    await redis.del(INVITE_LINK_PREFIX + token);

    const defaultChannel = await this.channelRepository.findOneOrFail({
      where: { guild },
      relations: ['guild'],
      order: { createdAt: 'ASC' }
    });

    const user = await this.userRepository.findOneOrFail(userId);

    this.socketService.addMember({ room: guild.id, member: user.toMember() });

    return {
      id: guild.id,
      name: guild.name,
      default_channel_id: defaultChannel.id,
      createdAt: guild?.createdAt.toString(),
      updatedAt: guild?.updatedAt.toString()
    };
  }

  async leaveGuild(userId: string, guildId: string): Promise<boolean> {
    const member = await this.memberRepository.findOneOrFail({ where: { guildId, userId } });
    await this.memberRepository.delete({ id: member.id });
    this.socketService.removeMember({ room: guildId, memberId: userId });
    return true;
  }
}
