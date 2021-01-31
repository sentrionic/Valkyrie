import { Body, Controller, Get, Param, Post, UseGuards } from '@nestjs/common';
import { ChannelService } from './channel.service';
import { AuthGuard } from '../guards/http/auth.guard';
import { GetUser } from '../config/user.decorator';
import { DMChannelResponse } from '../models/response/DMChannelResponse';
import { MemberGuard } from '../guards/http/member.guard';

@Controller('channels')
export class ChannelController {
  constructor(private readonly channelService: ChannelService) {
  }

  @Get('/:guildId')
  @UseGuards(MemberGuard)
  async getGuildChannels(
    @Param('guildId') guildId: string
  ) {
    return this.channelService.getGuildChannels(guildId);
  }

  @Post('/:guildId')
  @UseGuards(MemberGuard)
  async createChannel(
    @GetUser() userId: string,
    @Param('guildId') guildId: string,
    @Body('name') name: string,
    @Body('isPublic') isPublic: boolean = true,
    @Body('members') members: string[] = [],
  ): Promise<boolean> {
    return this.channelService.createChannel(
      guildId,
      name,
      isPublic,
      userId,
      members,
    );
  }

  @Post('/:guildId/dms')
  @UseGuards(AuthGuard)
  async getOrCreateChannel(
    @GetUser() userId: string,
    @Body('memberId') member: string,
    @Param('guildId') guildId: string,
  ): Promise<DMChannelResponse> {
    return this.channelService.getOrCreateChannel(guildId, member, userId);
  }
}
