import { Body, Controller, Param, Post, UseGuards } from '@nestjs/common';
import { ChannelService } from './channel.service';
import { AuthGuard } from '../config/auth.guard';
import { GetUser } from '../config/user.decorator';
import { DMChannelResponse } from '../models/response/DMChannelResponse';

@Controller('channel')
export class ChannelController {
  constructor(private readonly channelService: ChannelService) {}

  @Post("/:guildId/channels")
  @UseGuards(AuthGuard)
  async createChannel(
    @GetUser() userId: string,
    @Param('guildId') guildId: string,
    @Body('name') name: string,
    @Body('isPublic') isPublic?: boolean | null,
    @Body('members') members: string[] = [],
  ): Promise<ChannelResponse> {
    return this.channelService.createChannel(
      guildId,
      name,
      isPublic,
      userId,
      members,
    );
  }

  @Post("/:guildId/dms")
  @UseGuards(AuthGuard)
  async getOrCreateChannel(
    @GetUser() userId: string,
    @Body('memberId') member: string,
    @Param("guildId") guildId: string
  ): Promise<DMChannelResponse> {
    return this.channelService.getOrCreateChannel(guildId, member, userId);
  }
}
