import { CanActivate, ExecutionContext, Injectable } from '@nestjs/common';
import { Member } from '../../entities/member.entity';
import { Channel } from '../../entities/channel.entity';

@Injectable()
export class WsChannelGuard implements CanActivate {
  async canActivate(context: ExecutionContext): Promise<boolean> {
    const client = context.switchToWs().getClient();
    if (!client?.handshake?.session["userId"]) return false;

    const id = client?.handshake?.session["userId"];
    const channelID = context.getArgs()[1];

    if (!channelID) return false;

    const channel = await Channel.findOne({
      where: { id: channelID },
      relations: ['guild'],
    });

    const member = await Member.findOne({
      where: { guildId: channel.guild.id, userId: id },
    });

    return !!member;
  }
}