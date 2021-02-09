import { Module } from '@nestjs/common';
import { ChannelController } from './channel.controller';
import { ChannelService } from './channel.service';
import { User } from '../entities/user.entity';
import { Member } from '../entities/member.entity';
import { Guild } from '../entities/guild.entity';
import { Channel } from '../entities/channel.entity';
import { TypeOrmModule } from '@nestjs/typeorm';
import { SocketModule } from '../socket/socket.module';
import { PCMember } from '../entities/pcmember.entity';
import { DMMember } from '../entities/dmmember.entity';

@Module({
  imports: [
    TypeOrmModule.forFeature([Channel, Guild, Member, User, PCMember, DMMember]),
    SocketModule
  ],
  controllers: [ChannelController],
  providers: [ChannelService]
})
export class ChannelModule {
}
