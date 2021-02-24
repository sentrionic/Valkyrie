import { Injectable } from '@nestjs/common';
import { TypeOrmOptionsFactory, TypeOrmModuleOptions } from '@nestjs/typeorm';
import { PRODUCTION } from '../utils/constants';

@Injectable()
export class DatabaseConnectionService implements TypeOrmOptionsFactory {
  createTypeOrmOptions(): TypeOrmModuleOptions {
    return {
      name: 'default',
      type: 'postgres',
      url: process.env.DATABASE_URL,
      synchronize: !PRODUCTION,
      dropSchema: false,
      logging: true,
      entities: ['dist/**/*.entity.js'],
      migrations: ["dist/migrations/*.js"]
    };
  }
}
