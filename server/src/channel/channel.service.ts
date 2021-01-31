import { Injectable, UnauthorizedException } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { getManager, In, Repository } from 'typeorm';
import { Channel } from '../entities/channel.entity';
import { Member } from '../entities/member.entity';
import { User } from '../entities/user.entity';
import { Guild } from '../entities/guild.entity';
import { DMChannelResponse } from '../models/response/DMChannelResponse';
import { SocketService } from '../socket/socket.service';

@Injectable()
export class ChannelService {
  constructor(
    @InjectRepository(Channel) private channelRepository: Repository<Channel>,
    @InjectRepository(Guild) private guildRepository: Repository<Guild>,
    @InjectRepository(Member) private memberRepository: Repository<Member>,
    @InjectRepository(User) private userRepository: Repository<User>,
    private readonly socketService: SocketService
  ) {
  }

  async createChannel(
    guildId: string,
    name: string,
    isPublic: boolean = true,
    userId: string,
    members: string[],
  ): Promise<boolean> {
    const data = { name, public: isPublic };
    const memberPromise = await this.memberRepository.findOneOrFail({
      where: { guildId, userId },
    });
    const guildPromise = await this.guildRepository.findOneOrFail({
      where: { id: guildId },
    });

    const [member, guild] = await Promise.all([memberPromise, guildPromise]);

    if (!member.admin) throw new UnauthorizedException();

    let channel: Channel;

    await getManager().transaction(async (entityManager) => {
      channel = this.channelRepository.create(data);
      channel.guild = guild;
      await channel.save();

      // if (!isPublic) {
      //   members = members.filter((m) => m !== userId);
      //   members.push(userId);
      //   const pcmembers = members.map((m) => ({
      //     userId: m,
      //     channelId: channel.id,
      //   }));
      //   pcmembers.forEach((member) => {
      //     entityManager.insert(PCMember, {
      //       channelId: channel.id,
      //       userId: member.userId,
      //     });
      //   });
      // }

      await entityManager.save(channel);
    });

    const response: ChannelResponse = {
      id: channel!.id,
      name: channel!.name,
      public: true,
      createdAt: channel!.createdAt.toString(),
      updatedAt: channel!.updatedAt.toString(),
    }

    this.socketService.addChannel({ room: guildId, channel: response });
    return true;
  }


  async getGuildChannels(guildId: string): Promise<ChannelResponse[]> {
    const channels = await this.channelRepository.find({ where: { guild: guildId } });
    return channels.map(c => c.toJson());
  }

  async getOrCreateChannel(
    guildId: string,
    members: string,
    userId: string,
  ): Promise<DMChannelResponse> {
    const member = await this.memberRepository.findOne({
      where: { guildId, userId },
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
      [guildId],
    );

    if (data.length) {
      return data[0];
    }

    const users: User[] = await this.userRepository.find({
      where: { id: In(allMembers) },
    });

    const name: string = users.map((u) => u.id).join('/');

    const channelId = await getManager().transaction(async (entityManager) => {
      const channel = await this.channelRepository.create({
        name,
        isPublic: false,
        dm: true,
      });
      channel.guild = await this.guildRepository.findOneOrFail({
        where: { id: guildId },
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
      name,
    };
  }
}
