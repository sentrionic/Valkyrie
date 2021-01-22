import { Module } from '@nestjs/common';
import { GuildController } from './guild.controller';
import { GuildService } from './guild.service';
import { TypeOrmModule } from '@nestjs/typeorm';
import { User } from '../entities/user.entity';
import { Guild } from '../entities/guild.entity';
import { Member } from '../entities/member.entity';
import { Channel } from '../entities/channel.entity';

@Module({
  imports: [TypeOrmModule.forFeature([Guild, User, Member, Channel])],
  controllers: [GuildController],
  providers: [GuildService]
})
export class GuildModule {}
