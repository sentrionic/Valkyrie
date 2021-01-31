import { Global, Module } from '@nestjs/common';
import { SocketService } from './socket.service';
import { AppGateway } from './app.gateway';
import { TypeOrmModule } from '@nestjs/typeorm';
import { User } from '../entities/user.entity';

@Global()
@Module({
  imports: [
    TypeOrmModule.forFeature([User,])
  ],
  providers: [SocketService, AppGateway],
  exports: [SocketService],
})
export class SocketModule {}
