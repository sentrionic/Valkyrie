import {
  BadRequestException,
  Injectable,
  InternalServerErrorException,
  NotFoundException,
  UnauthorizedException
} from '@nestjs/common';
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
import { GuildInput } from '../models/input/GuildInput';
import { GuildMemberInput } from "../models/input/GuildMemberInput";
import { deleteFile, uploadAvatarToS3 } from '../utils/fileUtils';
import { BufferFile } from '../types/BufferFile';

@Injectable()
export class GuildService {
  constructor(
    @InjectRepository(Guild) private guildRepository: Repository<Guild>,
    @InjectRepository(Channel) private channelRepository: Repository<Channel>,
    @InjectRepository(User) private userRepository: Repository<User>,
    @InjectRepository(Member) private memberRepository: Repository<Member>,
    private socketService: SocketService,
  ) {
  }

  /**
   * Get the members of the given guild
   * @param guildId
   */
  async getGuildMembers(guildId: string): Promise<MemberResponse[]> {
    const manager = getManager();
    return await manager.query(
      `select distinct u.id,
                       u.username,
                       u.image,
                       u."isOnline",
                       u."createdAt",
                       u."updatedAt",
                       exists(select 1 from friends f where f.user = u.id) as "isFriend"
       from users as u
                join members m on u."id"::text = m."userId"
       where m."guildId" = $1
       order by u.username
      `,
      [guildId],
    );
  }

  /**
   * Get the guilds of the current user
   * @param userId
   */
  async getUserGuilds(userId: string): Promise<GuildResponse[]> {
    const manager = getManager();
    return await manager.query(
      `select distinct g."id",
                       g."name",
                       g."ownerId",
                       g."icon",
                       g."createdAt",
                       g."updatedAt",
                       (select c.id as "default_channel_id"
                        from channels c
                                 join guilds g on g.id = c."guildId"
                        where g.id = member."guildId"
                        order by c."createdAt"
                        limit 1)
       from guilds g
                join members as member
                     on g."id"::text = member."guildId"
       where member."userId" = $1
       order by g."createdAt";`,
      [userId],
    );
  }

  /**
   * Create a guild. Throws a BadRequest exception if the guild
   * limit is reached
   * @param name
   * @param userId
   */
  async createGuild(name: string, userId: string): Promise<GuildResponse> {

    await this.checkGuildLimit(userId);

    try {
      let guild: Guild | null = null;
      let channel: Channel | null = null;

      // Create guild with default channel and owner
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

      return this.toGuildResponse(guild, channel.id);
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

    await this.checkGuildLimit(userId);

    // Link includes the domain part
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

    // Check if already member
    const isMember = await this.memberRepository.findOne({ where: { userId, guildId } });

    if (isMember) {
      throw new BadRequestException();
    }

    await this.memberRepository.insert({ id: await idGenerator(), userId, guildId });

    await redis.del(INVITE_LINK_PREFIX + token);

    const defaultChannel = await this.channelRepository.findOneOrFail({
      where: { guild },
      relations: ['guild'],
      order: { createdAt: 'ASC' },
    });

    const user = await this.userRepository.findOneOrFail({ where: { id: userId }, relations: ['friends'] });

    this.socketService.addMember({ room: guild.id, member: user.toMember(userId) });

    return this.toGuildResponse(guild, defaultChannel.id);
  }

  async leaveGuild(userId: string, guildId: string): Promise<boolean> {
    const member = await this.memberRepository.findOneOrFail({ where: { guildId, userId } });
    const guild = await this.guildRepository.findOneOrFail({ where: { id: guildId } });

    if (guild.ownerId === userId) throw new BadRequestException('The owner cannot leave their server');

    await this.memberRepository.delete({ id: member.id });
    this.socketService.removeMember({ room: guildId, memberId: userId });
    return true;
  }

  async editGuild(userId: string, guildId: string, input: GuildInput, image: BufferFile): Promise<boolean> {

    const guild = await this.guildRepository.findOneOrFail({ where: { id: guildId } });

    if (guild.ownerId !== userId) {
      throw new UnauthorizedException();
    }

    let icon = input.image;
    if (image) {
      const directory = `valkyrie/guilds/${guildId}`;
      if (guild.icon) await deleteFile(guild.icon);
      icon = await uploadAvatarToS3(directory, image);
    }

    // Frontend sets the null as string
    if (icon === "null") icon = null;

    await this.guildRepository.update(guildId, {
      name: input.name ?? guild.name,
      icon
    });

    const updatedGuild = await this.guildRepository.findOneOrFail(guildId);

    this.socketService.editGuild(updatedGuild);

    return true;
  }

  async deleteGuild(userId: string, guildId: string): Promise<boolean> {
    const guild = await this.guildRepository.findOneOrFail({ where: { id: guildId } });

    if (guild.ownerId !== userId) {
      throw new UnauthorizedException();
    }

    let memberIds: any[];

    const manager = getManager();
    memberIds = await manager.query('delete from members where "guildId" = $1 returning members."userId";',
      [guildId]
    );
    await manager.query('delete from pcmembers where "channelId" = (select id from channels where "guildId" = $1);',
      [guildId]
    );

    await this.guildRepository.remove(guild);
    this.socketService.deleteGuild(memberIds[0], guildId);

    return true;
  }

  async changeMemberSettings(userId: string, guildId: string, input: GuildMemberInput): Promise<boolean> {

    const member = await this.memberRepository.findOne({
      where: {
        userId,
        guildId
      }
    });

    if (!member) throw new NotFoundException();

    const { nickname, color } = input;

    await this.memberRepository.update({ id: member.id }, {
      color,
      nickname
    });

    return true;
  }

  /**
   * Check if the user is in less than 100 servers.
   * Throws a BadRequestException if that is not the case.
   * @param userId
   */
  async checkGuildLimit(userId: string) {
    const count = await this.memberRepository.count({ userId });

    if (count >= 100) {
      throw new BadRequestException('Server Limit is 100');
    }
  }

  toGuildResponse(guild: Guild, defaultChannelId: string): GuildResponse {
    return {
      id: guild.id,
      name: guild.name,
      default_channel_id: defaultChannelId,
      ownerId: guild.ownerId,
      createdAt: guild?.createdAt.toString(),
      updatedAt: guild?.updatedAt.toString(),
    };
  }
}
