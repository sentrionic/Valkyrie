import { Body, Controller, Get, Param, Post, UseGuards } from '@nestjs/common';
import {
  ApiBody,
  ApiCookieAuth,
  ApiCreatedResponse,
  ApiOkResponse,
  ApiOperation,
  ApiUnauthorizedResponse
} from '@nestjs/swagger';
import { ChannelService } from './channel.service';
import { AuthGuard } from '../guards/http/auth.guard';
import { GetUser } from '../config/user.decorator';
import { DMChannelResponse } from '../models/response/DMChannelResponse';
import { ChannelResponse } from '../models/response/ChannelResponse';
import { MemberGuard } from '../guards/http/member.guard';
import { ChannelInput } from '../models/dto/ChannelInput';

@Controller('channels')
export class ChannelController {
  constructor(private readonly channelService: ChannelService) {
  }

  @Get('/:guildId')
  @UseGuards(MemberGuard)
  @ApiOperation({ summary: 'Get Guild Channels' })
  @ApiBody({ description: 'guildId', type: String })
  @ApiUnauthorizedResponse({ description: 'Invalid credentials' })
  @ApiCookieAuth()
  @ApiOkResponse({ type: [ChannelResponse] })
  async getGuildChannels(
    @Param('guildId') guildId: string
  ): Promise<ChannelResponse[]> {
    return this.channelService.getGuildChannels(guildId);
  }

  @Post('/:guildId')
  @UseGuards(MemberGuard)
  @ApiOperation({ summary: 'Create Guild Channels' })
  @ApiBody({ type: ChannelInput })
  @ApiUnauthorizedResponse({ description: 'Invalid credentials' })
  @ApiCookieAuth()
  @ApiCreatedResponse({ type: Boolean })
  async createChannel(
    @GetUser() userId: string,
    @Param('guildId') guildId: string,
    @Body() input: ChannelInput
  ): Promise<boolean> {
    return this.channelService.createChannel(
      guildId,
      userId,
      input
    );
  }

  @Post('/:guildId/dms')
  @UseGuards(AuthGuard)
  async getOrCreateChannel(
    @GetUser() userId: string,
    @Body('memberId') member: string,
    @Param('guildId') guildId: string
  ): Promise<DMChannelResponse> {
    return this.channelService.getOrCreateChannel(guildId, member, userId);
  }
}
