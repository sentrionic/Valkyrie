import { Injectable, UnauthorizedException } from '@nestjs/common';
import { Server, Socket } from 'socket.io';
import { MessageResponse } from '../models/response/MessageResponse';
import { MemberResponse } from '../models/response/MemberResponse';
import { getManager, Repository } from 'typeorm';
import { InjectRepository } from '@nestjs/typeorm';
import { User } from '../entities/user.entity';
import { ChannelResponse } from '../models/response/ChannelResponse';
import { Channel } from '../entities/channel.entity';
import { Member } from '../entities/member.entity';
import { PCMember } from '../entities/pcmember.entity';
import { WsException } from '@nestjs/websockets';

@Injectable()
export class SocketService {

  public socket: Server = null;

  constructor(
    @InjectRepository(User) private userRepository: Repository<User>,
    @InjectRepository(Channel) private channelRepository: Repository<Channel>,
    @InjectRepository(Member) private memberRepository: Repository<Member>,
    @InjectRepository(PCMember) private pcMemberRepository: Repository<PCMember>,
  ) {
  }

  async joinChannel(client: Socket, room: string) {
    const id: string = client.handshake.session['userId'];

    const channel = await this.channelRepository.findOneOrFail({
      where: { id: room },
      relations: ['guild']
    });

    await this.isChannelMember(channel, id);

    client.join(room);
  }

  sendMessage(
    message: { room: string; message: MessageResponse },
  ) {
    this.socket.to(message.room).emit('new_message', message.message);
  }

  editMessage(
    message: { room: string; message: MessageResponse },
  ) {
    this.socket.to(message.room).emit('edit_message', message.message);
  }

  deleteMessage(
    message: { room: string; message: MessageResponse },
  ) {
    this.socket.to(message.room).emit('delete_message', message.message);
  }

  addChannel(
    message: { room: string; channel: ChannelResponse },
  ) {
    this.socket.to(message.room).emit('add_channel', message.channel);
  }

  editChannel(
    message: { room: string; channel: ChannelResponse },
  ) {
    this.socket.to(message.room).emit('edit_channel', message.channel);
  }

  deleteChannel(
    message: { room: string, channelId: string },
  ) {
    this.socket.to(message.room).emit('delete_channel', message.channelId);
  }

  addMember(
    message: { room: string, member: MemberResponse },
  ) {
    this.socket.to(message.room).emit('add_member', message.member);
  }

  removeMember(
    message: { room: string, memberId: string },
  ) {
    this.socket.to(message.room).emit('remove_member', message.memberId);
  }

  async toggleOnlineStatus(client: Socket) {
    const id: string = client.handshake.session['userId'];
    await this.setOnlineStatus(id, true);
    const manager = getManager();
    const guilds: string[] = await manager.query(
      `select g.id from guilds g join members m on m."guildId" = g."id" where m."userId" = $1`,
      [id],
    );
    guilds.forEach(g => {
      const gId = g["id"];
      if (gId !== undefined) this.socket.to(gId).emit('toggle_online', id);
    });
  }

  async toggleOfflineStatus(client: Socket) {
    const id: string = client.handshake.session['userId'];
    await this.setOnlineStatus(id, false);
    const manager = getManager();
    const guilds: string[] = await manager.query(
      `select g.id from guilds g join members m on m."guildId" = g."id" where m."userId" = $1`,
      [id],
    );
    guilds.forEach(g => {
      const gId = g["id"];
      if (gId !== undefined) this.socket.to(gId).emit('toggle_offline', id);
    });
  }

  async setOnlineStatus(userId: string, isOnline: boolean): Promise<void> {
    await this.userRepository.update(userId, { isOnline });
  }

  addTyping(room: string, username: string) {
    this.socket.to(room).emit("addToTyping", username);
  }

  stopTyping(room: string, username: string) {
    this.socket.to(room).emit("removeFromTyping", username);
  }

  private async isChannelMember(channel: Channel, userId: string): Promise<boolean> {
    // Check if user has access to private channel
    if (!channel.isPublic) {
      const member = await this.pcMemberRepository.findOne({
        where: { channelId: channel.id, userId },
      });

      if (!member) {
        throw new WsException('Not Authorized');
      }
      // Check if user has access to the channel
    } else {
      const member = await this.memberRepository.findOneOrFail({
        where: { guildId: channel.guild.id, userId },
      });

      if (!member) {
        throw new WsException('Not Authorized');
      }
    }

    return true;
  }
}
