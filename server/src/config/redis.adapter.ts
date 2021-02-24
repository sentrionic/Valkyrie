import { IoAdapter } from '@nestjs/platform-socket.io';
import * as redisIoAdapter from 'socket.io-redis';
import { config } from 'dotenv';
import { redis } from './redis';

config();

export class RedisIoAdapter extends IoAdapter {
  createIOServer(port: number, options?: any): any {
    const server = super.createIOServer(port, options);
    const pubClient = redis;
    const subClient = redis.duplicate();
    const redisAdapter = redisIoAdapter({
      pubClient,
      subClient,
    });
    server.adapter(redisAdapter);
    return server;
  }
}
