import { CacheModule, Module, OnModuleInit } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import * as redisStore from 'cache-manager-redis-store';
import { DatabaseConnectionService } from './config/database';
import { UserModule } from './user/user.module';
import { GuildModule } from './guild/guild.module';
import { ChannelModule } from './channel/channel.module';
import { MessageModule } from './message/message.module';
import { SocketModule } from './socket/socket.module';
import { Connection } from 'typeorm';
import { PRODUCTION } from './utils/constants';

@Module({
  imports: [
    TypeOrmModule.forRootAsync({
      useClass: DatabaseConnectionService,
    }),
    CacheModule.register({
      store: redisStore,
      host: process.env.REDIS_HOST,
      port: process.env.REDIS_PORT,
      password: process.env.REDIS_PASSWORD,
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
export class AppModule implements OnModuleInit {
  constructor(private readonly connection: Connection) {}

  async onModuleInit(): Promise<void> {
    if (PRODUCTION) {
      await this.connection.runMigrations();
    }
  }
}
