import { Body, Controller, Delete, Get, Param, Post, UseGuards } from '@nestjs/common';
import { GuildService } from './guild.service';
import { AuthGuard } from '../guards/http/auth.guard';
import { GetUser } from '../config/user.decorator';
import { MemberGuard } from '../guards/http/member.guard';
import { MemberResponse } from '../models/response/MemberResponse';
import { YupValidationPipe } from '../utils/yupValidationPipe';
import { GuildSchema } from '../validation/guild.schema';
import { GuildInput } from '../models/input/GuildInput';
import { GuildResponse } from '../models/response/GuildResponse';
import { ApiBody, ApiCookieAuth, ApiOkResponse, ApiOperation, ApiUnauthorizedResponse } from '@nestjs/swagger';

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
    return this.guildService.getUserGuilds(userId);
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
    return this.guildService.createGuild(name, user);
  }

  @Get("/:guildId/invite")
  @UseGuards(MemberGuard)
  @ApiOperation({ summary: 'Create Invite Link' })
  @ApiUnauthorizedResponse({ description: 'Invalid credentials' })
  @ApiCookieAuth()
  @ApiBody({ type: String, description: "The guildId" })
  @ApiOkResponse({ type: String, description: "The invite link" })
  async generateTeamInvite(@Param('guildId') id: string): Promise<string> {
    return this.guildService.generateInviteLink(id);
  }

  @Post("/join")
  @UseGuards(AuthGuard)
  @ApiOperation({ summary: 'Join Guild' })
  @ApiUnauthorizedResponse({ description: 'Invalid credentials' })
  @ApiCookieAuth()
  @ApiBody({ type: String, description: "The invite link" })
  @ApiOkResponse({ type: GuildResponse })
  async joinTeam(
    @Body('link') link: string,
    @GetUser() user: string,
  ): Promise<GuildResponse> {
    return this.guildService.joinGuild(link, user);
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
    return this.guildService.leaveGuild(userId, guildId);
  }
}
