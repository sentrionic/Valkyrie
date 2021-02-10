import { CacheModule, Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import * as redisStore from 'cache-manager-redis-store';
import { DatabaseConnectionService } from './config/database';
import { UserModule } from './user/user.module';
import { GuildModule } from './guild/guild.module';
import { ChannelModule } from './channel/channel.module';
import { MessageModule } from './message/message.module';
import { SocketModule } from './socket/socket.module';

@Module({
  imports: [
    TypeOrmModule.forRootAsync({
      useClass: DatabaseConnectionService,
    }),
    CacheModule.register({
      store: redisStore,
      host: process.env.REDIS_URL_PUB_SUB,
      port: 6379,
    }),
    UserModule,
    GuildModule,
    ChannelModule,
    MessageModule,
    SocketModule,
  ],
  controllers: [],
  providers: [],
})
export class AppModule {}
