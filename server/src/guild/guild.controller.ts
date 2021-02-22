import {
  Body,
  Controller,
  Delete,
  Get,
  Param,
  Post,
  Put, Query, UploadedFile,
  UseGuards,
  UseInterceptors,
  ValidationPipe
} from '@nestjs/common';
import { GuildService } from './guild.service';
import { AuthGuard } from '../guards/http/auth.guard';
import { GetUser } from '../config/user.decorator';
import { MemberGuard } from '../guards/http/member.guard';
import { MemberResponse } from '../models/response/MemberResponse';
import { YupValidationPipe } from '../utils/yupValidationPipe';
import { GuildSchema, UpdateGuildSchema } from '../validation/guild.schema';
import { GuildInput } from '../models/input/GuildInput';
import { GuildResponse } from '../models/response/GuildResponse';
import {
  ApiBody,
  ApiConsumes,
  ApiCookieAuth,
  ApiOkResponse,
  ApiOperation, ApiTags,
  ApiUnauthorizedResponse
} from '@nestjs/swagger';
import { GuildMemberInput } from "../models/input/GuildMemberInput";
import { MemberSchema } from "../validation/member.schema";
import { FileInterceptor } from '@nestjs/platform-express';
import { BufferFile } from '../types/BufferFile';

@ApiTags('Guild Operation')
@Controller('guilds')
export class GuildController {
  constructor(private readonly guildService: GuildService) {
  }

  @Get("/:guildId/members")
  @UseGuards(MemberGuard)
  @ApiOperation({ summary: 'Get Guild Members' })
  @ApiUnauthorizedResponse({ description: 'Invalid credentials' })
  @ApiCookieAuth()
  @ApiOkResponse({ type: [MemberResponse] })
  async getGuildMembers(
    @Param('guildId') guildId: string,
    @GetUser() userId: string,
  ): Promise<MemberResponse[]> {
    return await this.guildService.getGuildMembers(guildId);
  }

  @Get()
  @UseGuards(AuthGuard)
  @ApiOperation({ summary: 'Get Users Guilds' })
  @ApiUnauthorizedResponse({ description: 'Invalid credentials' })
  @ApiCookieAuth()
  @ApiOkResponse({ type: [GuildResponse] })
  async getGuilds(
    @GetUser() userId: string
  ): Promise<GuildResponse[]> {
    return await this.guildService.getUserGuilds(userId);
  }

  @Post("/create")
  @UseGuards(AuthGuard)
  @ApiOperation({ summary: 'Create Guild' })
  @ApiUnauthorizedResponse({ description: 'Invalid credentials' })
  @ApiCookieAuth()
  @ApiBody({ type: GuildInput })
  @ApiOkResponse({ type: [GuildResponse] })
  async createGuild(
    @Body(new YupValidationPipe(GuildSchema)) input: GuildInput,
    @GetUser() user: string,
  ): Promise<GuildResponse> {
    const { name } = input;
    return await this.guildService.createGuild(name, user);
  }

  @Get("/:guildId/invite")
  @UseGuards(MemberGuard)
  @ApiOperation({ summary: 'Create Invite Link' })
  @ApiUnauthorizedResponse({ description: 'Invalid credentials' })
  @ApiCookieAuth()
  @ApiBody({ type: String, description: "The guildId" })
  @ApiOkResponse({ type: String, description: "The invite link" })
  async generateTeamInvite(
    @Param('guildId') id: string,
    @Query('isPermanent') isPermanent?: boolean
  ): Promise<string> {
    return await this.guildService.generateInviteLink(id, isPermanent);
  }

  @Delete("/:guildId/invite")
  @UseGuards(AuthGuard)
  @ApiOperation({ summary: 'Delete all permanent invite links' })
  @ApiUnauthorizedResponse({ description: 'Invalid credentials' })
  @ApiCookieAuth()
  @ApiBody({ type: String, description: "The guildId" })
  @ApiOkResponse({ type: Boolean })
  async deleteAllInvites(
    @Param('guildId') id: string,
    @GetUser() userId: string,
  ): Promise<boolean> {
    return await this.guildService.invalidateGuildInvites(id, userId);
  }

  @Post("/join")
  @UseGuards(AuthGuard)
  @ApiOperation({ summary: 'Join Guild' })
  @ApiUnauthorizedResponse({ description: 'Invalid credentials' })
  @ApiCookieAuth()
  @ApiBody({ type: String, description: "The invite link" })
  @ApiOkResponse({ type: GuildResponse })
  async joinGuild(
    @Body('link') link: string,
    @GetUser() user: string,
  ): Promise<GuildResponse> {
    return await this.guildService.joinGuild(link, user);
  }

  @Get("/:guildId/member")
  @UseGuards(MemberGuard)
  @ApiOperation({ summary: 'Get Member Settings' })
  @ApiUnauthorizedResponse({ description: 'Invalid credentials' })
  @ApiCookieAuth()
  async getMemberSettings(
    @GetUser() user: string,
    @Param('guildId') guildId: string,
  ): Promise<GuildMemberInput> {
    return await this.guildService.getMemberSettings(user, guildId);
  }

  @Put("/:guildId/member")
  @UseGuards(MemberGuard)
  @ApiBody({ type: GuildMemberInput })
  @ApiOperation({ summary: 'Edit Member Settings' })
  @ApiUnauthorizedResponse({ description: 'Invalid credentials' })
  @ApiCookieAuth()
  async editMember(
    @GetUser() user: string,
    @Param('guildId') guildId: string,
    @Body(
      new YupValidationPipe(MemberSchema),
      new ValidationPipe({ transform: true })
    ) input: GuildMemberInput,
  ): Promise<boolean> {
    return await this.guildService.changeMemberSettings(user, guildId, input);
  }

  @Delete("/:guildId")
  @UseGuards(AuthGuard)
  @ApiOperation({ summary: 'Leave Guild' })
  @ApiUnauthorizedResponse({ description: 'Invalid credentials' })
  @ApiCookieAuth()
  @ApiOkResponse({ type: Boolean })
  async leaveGuild(
    @GetUser() userId: string,
    @Param("guildId") guildId: string
  ): Promise<boolean> {
    return await this.guildService.leaveGuild(userId, guildId);
  }

  @Put("/:guildId")
  @UseGuards(AuthGuard)
  @UseInterceptors(FileInterceptor('image'))
  @ApiCookieAuth()
  @ApiOperation({ summary: 'Edit Guild' })
  @ApiOkResponse({ description: 'Edit Success', type: Boolean })
  @ApiUnauthorizedResponse()
  @ApiBody({ type: GuildInput })
  @ApiConsumes('multipart/form-data')
  async editGuild(
    @GetUser() user: string,
    @Param('guildId') guildId: string,
    @Body(
      new YupValidationPipe(UpdateGuildSchema),
      new ValidationPipe({ transform: true })
    ) input: GuildInput,
    @UploadedFile() image?: BufferFile,
  ): Promise<boolean> {
    return await this.guildService.editGuild(user, guildId, input, image);
  }

  @Delete("/:guildId/delete")
  @UseGuards(AuthGuard)
  @ApiCookieAuth()
  @ApiOperation({ summary: 'Delete Guild' })
  @ApiOkResponse({ description: 'Delete Success', type: Boolean })
  @ApiUnauthorizedResponse()
  async deleteGuild(
    @GetUser() userId: string,
    @Param('guildId') guildId: string,
  ): Promise<boolean> {
    return await this.guildService.deleteGuild(userId, guildId);
  }

  @Get("/:guildId/bans")
  @UseGuards(AuthGuard)
  @ApiCookieAuth()
  @ApiOperation({ summary: 'Get Guild\'s ban list' })
  @ApiOkResponse({ description: 'List of users', type: [MemberResponse] })
  @ApiUnauthorizedResponse()
  async getBannedUsers(
    @GetUser() userId: string,
    @Param('guildId') guildId: string,
  ): Promise<MemberResponse[]> {
    return await this.guildService.getBannedUsers(userId, guildId);
  }

  @Post("/:guildId/bans")
  @UseGuards(AuthGuard)
  @ApiCookieAuth()
  @ApiOperation({ summary: 'Ban a user' })
  @ApiOkResponse({ description: 'Successfully banned', type: Boolean })
  @ApiBody({ type: String, description: "MemberId" })
  @ApiUnauthorizedResponse()
  async banUser(
    @GetUser() userId: string,
    @Param('guildId') guildId: string,
    @Body('memberId') memberId: string
  ): Promise<boolean> {
    return await this.guildService.banMember(userId, guildId, memberId);
  }

  @Post("/:guildId/kick")
  @UseGuards(AuthGuard)
  @ApiCookieAuth()
  @ApiOperation({ summary: 'Kick a user' })
  @ApiOkResponse({ description: 'Successfully kicked', type: Boolean })
  @ApiBody({ type: String, description: "MemberId" })
  @ApiUnauthorizedResponse()
  async kickUser(
    @GetUser() userId: string,
    @Param('guildId') guildId: string,
    @Body('memberId') memberId: string
  ): Promise<boolean> {
    return await this.guildService.kickMember(userId, guildId, memberId);
  }

  @Delete("/:guildId/bans")
  @UseGuards(AuthGuard)
  @ApiCookieAuth()
  @ApiOperation({ summary: 'Unban a user' })
  @ApiOkResponse({ description: 'Successfully unbanned', type: Boolean })
  @ApiBody({ type: String, description: "MemberId" })
  @ApiUnauthorizedResponse()
  async unbanUser(
    @GetUser() userId: string,
    @Param('guildId') guildId: string,
    @Body('memberId') memberId: string
  ): Promise<boolean> {
    return await this.guildService.unbanUser(userId, guildId, memberId);
  }
}
