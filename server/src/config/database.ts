import { Injectable } from '@nestjs/common';
import { TypeOrmOptionsFactory, TypeOrmModuleOptions } from '@nestjs/typeorm';

@Injectable()
export class DatabaseConnectionService implements TypeOrmOptionsFactory {
  createTypeOrmOptions(): TypeOrmModuleOptions {
    return {
      name: 'default',
      type: 'postgres',
      url: process.env.DATABASE_URL,
      synchronize: true,
      dropSchema: false,
      logging: true,
      entities: ['dist/**/*.entity.js'],
    };
  }
}
