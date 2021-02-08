import { Injectable, NotFoundException, UnauthorizedException } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { getManager, In, Repository } from 'typeorm';
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

@Injectable()
export class ChannelService {
  constructor(
    @InjectRepository(Channel) private channelRepository: Repository<Channel>,
    @InjectRepository(Guild) private guildRepository: Repository<Guild>,
    @InjectRepository(Member) private memberRepository: Repository<Member>,
    @InjectRepository(User) private userRepository: Repository<User>,
    @InjectRepository(PCMember) private pcMemberRepository: Repository<PCMember>,
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
    const memberPromise = await this.memberRepository.findOneOrFail({
      where: { guildId, userId }
    });
    const guildPromise = await this.guildRepository.findOneOrFail({
      where: { id: guildId }
    });

    const [member, guild] = await Promise.all([memberPromise, guildPromise]);

    if (!member.admin) throw new UnauthorizedException();

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
    guildId: string,
    members: string,
    userId: string
  ): Promise<DMChannelResponse> {
    const member = await this.memberRepository.findOne({
      where: { guildId, userId }
    });

    if (!member) {
      throw new Error('Not Authorized');
    }

    const allMembers = [...members, userId];

    // create string containing all member ids and seperate them with a comma
    let array = '';
    allMembers.forEach((member, index, arr) => {
      array += `'${member}'`;
      if (index < arr.length - 1) array += ',';
    });

    // check if dm channel already exists with these members
    const data = await getManager().query(
      `
        select c.id, c.name 
        from channels as c, pcmembers pc 
        where pc."channelId"::text = c."id"::text and c.dm = true and c.public = false and c."guildId"::text = $1
        group by c."id", c."name"  
        having array_agg(pc."userId"::text) @> Array[${array}]
        and count(pc."userId") = ${allMembers.length};
        `,
      [guildId]
    );

    if (data.length) {
      return data[0];
    }

    const users: User[] = await this.userRepository.find({
      where: { id: In(allMembers) }
    });

    const name: string = users.map((u) => u.id).join('/');

    const channelId = await getManager().transaction(async (entityManager) => {
      const channel = await this.channelRepository.create({
        name,
        isPublic: false,
        dm: true
      });
      channel.guild = await this.guildRepository.findOneOrFail({
        where: { id: guildId }
      });

      await channel.save();
      await entityManager.save(channel);

      const channelId = channel.id;
      // const pcmembers = allMembers.map((m) => ({ userId: m, channelId }));
      // pcmembers.forEach((member) => {
      //   entityManager.insert(PCMember, {
      //     channelId,
      //     userId: member.userId,
      //   });
      // });

      return channelId;
    });

    return {
      id: channelId,
      name
    };
  }

  async editChannel(userId: string, channelId: string, input: ChannelInput): Promise<boolean> {
    const channel = await this.channelRepository.findOneOrFail({
      where: { id: channelId },
      relations: ['guild'],
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
      relations: ['guild', 'members'],
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
      relations: ['guild'],
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
}
