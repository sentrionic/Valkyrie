import { Body, Controller, Get, Param, Post, UseGuards } from '@nestjs/common';
import { GuildService } from './guild.service';
import { AuthGuard } from '../config/auth.guard';
import { GetUser } from '../config/user.decorator';
import { Guild } from '../entities/guild.entity';
import { MemberGuard } from '../config/member.guard';
import { MemberResponse } from '../models/response/MemberResponse';

@Controller('guilds')
export class GuildController {
  constructor(private readonly guildService: GuildService) {
  }

  @Get("/:id/members")
  @UseGuards(MemberGuard)
  async getGuildMembers(
    @Param('id') guildId: string,
    @GetUser() userId: string,
  ): Promise<MemberResponse[]> {
    return await this.guildService.getGuildMembers(guildId);
  }

  // @ResolveField(() => [MemberResponse!]!)
  // @UseGuards(TeamGuard)
  // async directMessageMembers(
  //   @Parent() team: Team,
  //   @GetUser() user: User,
  // ): Promise<User[]> {
  //   return this.teamService.getDirectMessageMembers(team.id, user.id);
  // }

  @Post("/create")
  @UseGuards(AuthGuard)
  async createGuild(
    @Body('name') name: string,
    @GetUser() user: string,
  ): Promise<Guild> {
    return this.guildService.createGuild(name, user);
  }

  @Post("/:id/invite")
  @UseGuards(MemberGuard)
  async generateTeamInvite(@Param('id') id: string): Promise<string> {
    return this.guildService.generateInviteLink(id);
  }

  @Post("/join")
  @UseGuards(AuthGuard)
  async joinTeam(
    @Body('token') token: string,
    @GetUser() user: string,
  ): Promise<Guild> {
    return this.guildService.joinGuild(token, user);
  }
}
