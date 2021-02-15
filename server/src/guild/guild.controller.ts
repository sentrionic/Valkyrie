import {
  Body,
  Controller,
  Delete,
  Get,
  Param,
  Post,
  Put, UploadedFile,
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
  ApiOperation,
  ApiUnauthorizedResponse
} from '@nestjs/swagger';
import { GuildMemberInput } from "../models/input/GuildMemberInput";
import { MemberSchema } from "../validation/member.schema";
import { FileInterceptor } from '@nestjs/platform-express';
import { BufferFile } from '../types/BufferFile';

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
  async joinGuild(
    @Body('link') link: string,
    @GetUser() user: string,
  ): Promise<GuildResponse> {
    return this.guildService.joinGuild(link, user);
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
    return this.guildService.changeMemberSettings(user, guildId, input);
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
    return this.guildService.editGuild(user, guildId, input, image);
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
    return this.guildService.deleteGuild(userId, guildId);
  }
}
