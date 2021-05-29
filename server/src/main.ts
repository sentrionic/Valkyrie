import { NestFactory } from '@nestjs/core';
import { NestExpressApplication } from '@nestjs/platform-express';
import { DocumentBuilder, SwaggerModule } from '@nestjs/swagger';
import * as helmet from 'helmet';
import { config } from 'dotenv';
import * as rateLimit from 'express-rate-limit';
import { AppModule } from './app.module';
import { COOKIE_NAME } from './utils/constants';
import { sessionMiddleware } from './config/sessionmiddleware';
import { RedisIoAdapter } from './config/redis.adapter';
import { redis } from './config/redis';
import { join } from 'path';

const RedisStore = require('rate-limit-redis');

config();

async function bootstrap() {
  const app = await NestFactory.create<NestExpressApplication>(AppModule);
  app.useWebSocketAdapter(new RedisIoAdapter(app));
  app.useStaticAssets(join(__dirname, '..', 'static'), { prefix: '/ws' });
  app.setGlobalPrefix('api');
  app.set('trust proxy', 1);
  app.use(helmet());
  app.enableCors({
    origin: process.env.CORS_ORIGIN,
    credentials: true,
  });

  app.use(sessionMiddleware);
  app.use(
    rateLimit({
      store: new RedisStore({
        client: redis,
      }),
      windowMs: 60 * 1000, // 1 minutes
      max: 100, // limit each IP to 100 requests per windowMs
    }),
  );

  const options = new DocumentBuilder()
    .setTitle('Valkyrie API')
    .setDescription(
      'Valkyrie REST API Specs. ' +
        'For the Websocket Endpoints and events go <a href="/ws">here</a>. ' +
        'Both services use <a href="https://github.com/expressjs/session">Express Sessions</a> ' +
        'for authentication.',
    )
    .setVersion('1.0.0')
    .addCookieAuth(COOKIE_NAME, {
      type: 'http',
    })
    .build();

  const document = SwaggerModule.createDocument(app, options);
  SwaggerModule.setup('/', app, document);
  await app.listen(process.env.PORT || 4000);
}

bootstrap();
