import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { DatabaseConnectionService } from './config/database';
import { UserModule } from './user/user.module';
import { GuildModule } from './guild/guild.module';
import { ChannelModule } from './channel/channel.module';
import { MessageModule } from './message/message.module';

@Module({
  imports: [
    TypeOrmModule.forRootAsync({
      useClass: DatabaseConnectionService,
    }),
    UserModule,
    GuildModule,
    ChannelModule,
    MessageModule,
  ],
  controllers: [],
  providers: [],
})
export class AppModule {}
