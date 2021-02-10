import { BadRequestException, Injectable, NotFoundException, UnauthorizedException } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { getManager, Repository } from 'typeorm';
import { Channel } from '../entities/channel.entity';
import { Member } from '../entities/member.entity';
import { User } from '../entities/user.entity';
import { Guild } from '../entities/guild.entity';
import { DMChannelResponse } from '../models/response/DMChannelResponse';
import { SocketService } from '../socket/socket.service';
import { ChannelResponse } from '../models/response/ChannelResponse';
import { ChannelInput } from '../models/dto/ChannelInput';
import { PCMember } from '../entities/pcmember.entity';
import { idGenerator } from '../utils/idGenerator';
import { DMMember } from '../entities/dmmember.entity';

@Injectable()
export class ChannelService {
  constructor(
    @InjectRepository(Channel) private channelRepository: Repository<Channel>,
    @InjectRepository(Guild) private guildRepository: Repository<Guild>,
    @InjectRepository(Member) private memberRepository: Repository<Member>,
    @InjectRepository(User) private userRepository: Repository<User>,
    @InjectRepository(PCMember) private pcMemberRepository: Repository<PCMember>,
    @InjectRepository(DMMember) private dmMemberRepository: Repository<DMMember>,
    private readonly socketService: SocketService
  ) {
  }

  async createChannel(
    guildId: string,
    userId: string,
    input: ChannelInput
  ): Promise<boolean> {

    const { name, isPublic } = input;
    let { members } = input;

    const data = { name: name.trim(), public: isPublic };

    const guild = await this.guildRepository.findOneOrFail({
      where: { id: guildId }
    });

    if (guild.ownerId !== userId) throw new UnauthorizedException();

    const count = await this.channelRepository.count({ guild });

    if (count >= 50) {
      throw new BadRequestException('Channel Limit is 50');
    }

    let channel: Channel;

    await getManager().transaction(async (entityManager) => {
      channel = this.channelRepository.create(data);
      channel.guild = guild;
      await channel.save();

      if (!isPublic) {
        members = members.filter((m) => m !== userId);
        members.push(userId);
        const pcmembers = members.map((m) => ({
          userId: m,
          channelId: channel.id
        }));

        pcmembers.forEach((member) => {
          entityManager.insert(PCMember, {
            id: idGenerator(),
            channelId: channel.id,
            userId: member.userId
          });
        });

        channel.isPublic = false;
      }

      await entityManager.save(channel);
    });

    const response: ChannelResponse = {
      id: channel!.id,
      name: channel!.name,
      isPublic: channel!.isPublic,
      createdAt: channel!.createdAt.toString(),
      updatedAt: channel!.updatedAt.toString()
    };

    this.socketService.addChannel({ room: guildId, channel: response });
    return true;
  }


  async getGuildChannels(guildId: string, userId: string): Promise<ChannelResponse[]> {
    const manager = getManager();
    return await manager.query(
      `
          select distinct on (c.id, c."createdAt") c.id, c.name, c."isPublic", c."createdAt", c."updatedAt"
          from channels as c
                   left outer join pcmembers as pc
                                   on c."id"::text = pc."channelId"::text
          where c."guildId"::text = $1
            and (c."isPublic" = true or pc."userId"::text = $2)
          order by c."createdAt"
      `,
      [guildId, userId]
    );
  }

  async getOrCreateChannel(
    userId: string,
    memberId: string
  ): Promise<DMChannelResponse> {
    const member = await this.userRepository.findOne({
      where: { id: memberId }
    });

    if (!member) {
      throw new NotFoundException();
    }

    // check if dm channel already exists with these members
    const data = await getManager().query(
      `
        select c.id
        from channels as c, dm_members dm 
        where dm."channelId" = c."id" and c.dm = true and c."isPublic" = false
        group by c."id"
        having array_agg(dm."userId"::text) @> Array['${memberId}', '${userId}']
        and count(dm."userId") = 2;
        `
    );

    if (data.length) {
      this.setDirectMessageStatus(data[0].id, userId, true);
      return {
        id: data[0].id,
        user: member.toMember()
      };
    }

    const channelId = await getManager().transaction(async (entityManager) => {
      const channel = await this.channelRepository.create({
        name: idGenerator(),
        isPublic: false,
        dm: true
      });
      await channel.save();
      await entityManager.save(channel);

      const channelId = channel.id;
      const allMembers = [memberId, userId];
      const dmMembers = allMembers.map((m) => ({ userId: m, channelId }));
      dmMembers.forEach((member) => {
        entityManager.insert(DMMember, {
          id: idGenerator(),
          channelId,
          userId: member.userId,
          isOpen: member.userId === userId
        });
      });

      return channelId;
    });

    return {
      id: channelId,
      user: member.toMember()
    };
  }

  async getDirectMessageChannels(userId: string): Promise<DMChannelResponse[]> {
    const manager = getManager();
    const result = await manager.query(
      `
          select dm."channelId", u.username, u.image, u.id, u."isOnline", u."createdAt", u."updatedAt"
          from users u
                   join dm_members dm on dm."userId" = u.id
          where u.id != $1
            and dm."channelId" in (
              select distinct c.id
              from channels as c
                       left outer join dm_members as dm
                                       on c."id" = dm."channelId"
                       join users u on dm."userId" = u.id
              where c."isPublic" = false
                and c.dm = true
                and dm."isOpen" = true
                and dm."userId" = $1
          )
      `,
      [userId]
    );

    const dms: DMChannelResponse[] = [];
    result.map(r => dms.push({
      id: r.channelId,
      user: {
        id: r.id,
        username: r.username,
        image: r.image,
        isOnline: r.isOnline,
        createdAt: r.createdAt,
        updatedAt: r.updatedAt
      }
    }));
    return dms;
  }

  async editChannel(userId: string, channelId: string, input: ChannelInput): Promise<boolean> {
    const channel = await this.channelRepository.findOneOrFail({
      where: { id: channelId },
      relations: ['guild']
    });

    if (!channel) {
      throw new NotFoundException();
    }

    if (channel.guild.ownerId !== userId) {
      throw new UnauthorizedException();
    }

    const { name, isPublic } = input;
    let { members } = input;

    // Used to be private and now is public
    if (isPublic && !channel.isPublic) {
      await getManager().query(
        'delete from pcmembers where "channelId" = $1;',
        [channelId]
      );
    }

    await this.channelRepository.update(channelId, {
      name: name ?? channel.name,
      isPublic: isPublic ?? channel.isPublic
    });

    // Member Changes
    if (!isPublic && members) {

      await getManager().transaction(async (entityManager) => {
        members = members.filter((m) => m !== userId);
        members.push(userId);

        const current = await this.pcMemberRepository.find({ where: { channelId } });
        const newMembers = members.filter((m => !current.map(c => c.userId).includes(m)));
        const remove = current.filter((c => !members.map(m => m).includes(c.userId)));

        const pcmembers = newMembers.map((m) => ({
          userId: m,
          channelId: channel.id
        }));

        pcmembers.forEach((member) => {
          entityManager.insert(PCMember, {
            id: idGenerator(),
            channelId: channel.id,
            userId: member.userId
          });
        });

        if (remove.length > 0) {
          await entityManager.query(
            'delete from pcmembers where "userId" IN ($1) and "channelId" = $2;',
            [...remove.map(r => r.userId), channelId]
          );
        }
      });

    }

    const updatedChannel = await this.channelRepository.findOne({
      where: { id: channelId }
    });

    const response: ChannelResponse = {
      id: updatedChannel!.id,
      name: updatedChannel!.name,
      isPublic: updatedChannel!.isPublic,
      createdAt: updatedChannel!.createdAt.toString(),
      updatedAt: updatedChannel!.updatedAt.toString()
    };

    this.socketService.editChannel({ room: channel.guild.id, channel: response });

    return true;
  }

  async deleteChannel(userId: string, channelId: string): Promise<boolean> {
    const channel = await this.channelRepository.findOneOrFail({
      where: { id: channelId },
      relations: ['guild', 'members']
    });

    if (!channel) {
      throw new NotFoundException();
    }

    if (channel.guild.ownerId !== userId) {
      throw new UnauthorizedException();
    }

    if (!channel.isPublic) {
      await getManager().query(
        'delete from pcmembers where "channelId" = $1;',
        [channelId]
      );
    }

    await this.channelRepository.remove(channel);

    this.socketService.deleteChannel({ room: channel.guild.id, channelId });

    return true;
  }

  async getPrivateChannelMembers(userId: string, channelId: string): Promise<string[]> {
    const channel = await this.channelRepository.findOneOrFail({
      where: { id: channelId },
      relations: ['guild']
    });

    if (!channel) {
      throw new NotFoundException();
    }

    if (channel.guild.ownerId !== userId) {
      throw new UnauthorizedException();
    }

    if (channel.isPublic) return [];

    const ids = await getManager().query(
      `
          select pc."userId"
          from pcmembers pc
                   join channels c on pc."channelId" = c.id
          where c.id = $1
      `,
      [channelId]
    );

    if (ids.length === 0) return [];

    return ids.map(i => i.userId);
  }

  async setDirectMessageStatus(channelId: string, userId: string, isOpen: boolean): Promise<boolean> {
    const channel = await this.dmMemberRepository.findOneOrFail({
      where: { channelId, userId }
    });

    if (!channel) throw new NotFoundException();

    await this.dmMemberRepository.update({ id: channel.id }, {
      isOpen
    });

    return true;
  }
}
