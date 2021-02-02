import { CanActivate, ExecutionContext, Injectable } from '@nestjs/common';
import { Member } from '../../entities/member.entity';
import { Channel } from '../../entities/channel.entity';
import e from 'express';

@Injectable()
export class ChannelGuard implements CanActivate {
  async canActivate(context: ExecutionContext): Promise<boolean> {
    const request: e.Request = context.switchToHttp().getRequest();
    if (!(request?.session["userId"])) return false;

    const { channelId } = request.params;
    const id = request.session["userId"];

    const channel = await Channel.findOneOrFail({
      where: { id: channelId },
      relations: ['guild'],
    });

    const member = await Member.findOneOrFail({
      where: { guildId: channel.guild.id, userId: id },
    });

    return !!member;
  }
}
