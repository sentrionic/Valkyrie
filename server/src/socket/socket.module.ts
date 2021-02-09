import { Global, Module } from '@nestjs/common';
import { SocketService } from './socket.service';
import { AppGateway } from './app.gateway';
import { TypeOrmModule } from '@nestjs/typeorm';
import { User } from '../entities/user.entity';
import { Channel } from '../entities/channel.entity';
import { Member } from '../entities/member.entity';
import { PCMember } from '../entities/pcmember.entity';
import { DMMember } from '../entities/dmmember.entity';

@Global()
@Module({
  imports: [
    TypeOrmModule.forFeature([User, Channel, Member, PCMember, DMMember])
  ],
  providers: [SocketService, AppGateway],
  exports: [SocketService],
})
export class SocketModule {}
