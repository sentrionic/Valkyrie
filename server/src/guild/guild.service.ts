import { Injectable, InternalServerErrorException, NotFoundException } from '@nestjs/common';
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

@Injectable()
export class GuildService {
  constructor(
    @InjectRepository(Guild) private guildRepository: Repository<Guild>,
    @InjectRepository(Channel) private channelRepository: Repository<Channel>,
    @InjectRepository(User) private userRepository: Repository<User>,
    @InjectRepository(Member) private memberRepository: Repository<Member>,
  ) {
  }

  async getGuildMembers(guildId: string): Promise<MemberResponse[]> {
    const manager = getManager();
    return await manager.query(
      `select *
       from users as u
                join members m on u."id"::text = m."userId"
       where m."guildId" = $1`,
      [guildId],
    );
  }

  async getUserGuilds(userId: string): Promise<Guild[]> {
    const manager = getManager();
    return await manager.query(
      `select *
       from guild
                join members as member on guild."id"::text = member."guildId"
       where member."userId" = $1
       order by guild."createdAt" ASC;`,
      [userId],
    );
  }

  async getDirectMessageMembers(
    guildId: string,
    userId: string,
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
      [userId, guildId],
    );
  }

  async isAdmin(guildId: string, user: User): Promise<boolean> {
    return (
      await this.memberRepository.findOne({
        where: { guildId, userId: user.id },
      })
    ).admin;
  }

  async createGuild(name: string, userId: string): Promise<Guild> {
    try {
      let guild = null;

      await getManager().transaction(async (entityManager) => {
        guild = this.guildRepository.create();
        const channel = this.channelRepository.create({ name: 'general' });

        guild.name = name;
        await guild.save();
        await entityManager.save(guild);

        channel.guild = guild;
        await channel.save();
        await entityManager.save(channel);

        await entityManager.insert(Member, {
          guildId: guild.id,
          userId,
          admin: true,
        });
      });

      return guild;
    } catch (err) {
      throw new InternalServerErrorException(err);
    }
  }

  async generateInviteLink(guildId: string): Promise<string> {
    const token = nanoid(8);
    await redis.set(INVITE_LINK_PREFIX + token, guildId, 'ex', 60 * 60 * 24); // 1 day expiration
    return `${process.env.BACKEND_HOST}/${token}`;
  }

  async joinGuild(token: string, userId: string): Promise<Guild> {
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

    try {
      await this.memberRepository.insert({ userId, guildId });
    } catch (err) {
      if (err.code === '23505') {
        throw new InternalServerErrorException();
      }
    }

    await redis.del(INVITE_LINK_PREFIX + token);

    return guild;
  }
}
