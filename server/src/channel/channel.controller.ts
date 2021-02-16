import { Body, Controller, Delete, Get, Param, Post, Put, UseGuards, ValidationPipe } from '@nestjs/common';
import {
  ApiBody,
  ApiCookieAuth,
  ApiCreatedResponse,
  ApiOkResponse,
  ApiOperation, ApiTags,
  ApiUnauthorizedResponse
} from '@nestjs/swagger';
import { ChannelService } from './channel.service';
import { AuthGuard } from '../guards/http/auth.guard';
import { GetUser } from '../config/user.decorator';
import { DMChannelResponse } from '../models/response/DMChannelResponse';
import { ChannelResponse } from '../models/response/ChannelResponse';
import { MemberGuard } from '../guards/http/member.guard';
import { ChannelInput } from '../models/input/ChannelInput';
import { YupValidationPipe } from '../utils/yupValidationPipe';
import { ChannelSchema } from '../validation/channel.schema';

@ApiTags('Channel Operation')
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
    @Param('guildId') guildId: string,
    @GetUser() userId: string
  ): Promise<ChannelResponse[]> {
    return this.channelService.getGuildChannels(guildId, userId);
  }

  @Get('/:channelId/members')
  @UseGuards(AuthGuard)
  @ApiOperation({ summary: 'Get Private Guild Members' })
  @ApiBody({ description: 'channelId', type: String })
  @ApiUnauthorizedResponse({ description: 'Invalid credentials' })
  @ApiCookieAuth()
  @ApiOkResponse({ description: 'Member Ids', type: [String] })
  async getPrivateChannelMembers(
    @Param('channelId') channelId: string,
    @GetUser() userId: string
  ): Promise<string[]> {
    return this.channelService.getPrivateChannelMembers(userId, channelId);
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
    @Body(
      new YupValidationPipe(ChannelSchema),
      new ValidationPipe({ transform: true })
    ) input: ChannelInput
  ): Promise<boolean> {
    return this.channelService.createChannel(
      guildId,
      userId,
      input
    );
  }

  @Get('/me/dm')
  @UseGuards(AuthGuard)
  @ApiOperation({ summary: 'Get User\'s DMs' })
  @ApiUnauthorizedResponse({ description: 'Invalid credentials' })
  @ApiCookieAuth()
  @ApiOkResponse({ type: [DMChannelResponse] })
  async getDirectMessageChannels(
    @GetUser() userId: string
  ): Promise<DMChannelResponse[]> {
    return this.channelService.getDirectMessageChannels(userId);
  }

  @Post(':memberId/dm')
  @UseGuards(AuthGuard)
  @ApiOperation({ summary: 'Start or get DMs with the given user' })
  @ApiUnauthorizedResponse({ description: 'Invalid credentials' })
  @ApiCookieAuth()
  @ApiOkResponse({ type: DMChannelResponse })
  async getOrCreateChannel(
    @GetUser() userId: string,
    @Param('memberId') memberId: string
  ): Promise<DMChannelResponse> {
    return this.channelService.getOrCreateChannel(userId, memberId);
  }

  @Put("/:guildId/:channelId")
  @UseGuards(AuthGuard)
  @ApiCookieAuth()
  @ApiOperation({ summary: 'Edit Channel' })
  @ApiOkResponse({ description: 'Edit Success', type: Boolean })
  @ApiUnauthorizedResponse()
  @ApiBody({ type: ChannelInput })
  async editChannel(
    @GetUser() user: string,
    @Param('channelId') channelId: string,
    @Body(
      new YupValidationPipe(ChannelSchema),
      new ValidationPipe({ transform: true })
    ) input: ChannelInput,
  ): Promise<boolean> {
    return this.channelService.editChannel(user, channelId, input);
  }

  @Delete('/:channelId/dm')
  @UseGuards(AuthGuard)
  @ApiCookieAuth()
  @ApiOperation({ summary: 'Close the DM' })
  @ApiOkResponse({ description: 'Close Success', type: Boolean })
  @ApiUnauthorizedResponse()
  async closeDirectMessage(
    @GetUser() userId: string,
    @Param('channelId') channelId: string,
  ): Promise<boolean> {
    return this.channelService.setDirectMessageStatus(channelId, userId, false);
  }

  @Delete("/:guildId/:channelId")
  @UseGuards(MemberGuard)
  @ApiCookieAuth()
  @ApiOperation({ summary: 'Delete Channel' })
  @ApiOkResponse({ description: 'Delete Success', type: Boolean })
  @ApiUnauthorizedResponse()
  async deleteChannel(
    @GetUser() userId: string,
    @Param('channelId') channelId: string,
  ): Promise<boolean> {
    return this.channelService.deleteChannel(userId, channelId);
  }
}
