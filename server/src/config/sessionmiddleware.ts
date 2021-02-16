import * as session from 'express-session';
import { COOKIE_NAME } from '../utils/constants';
import { redis } from './redis';
import * as connectRedis from 'connect-redis';
import { config } from 'dotenv';
config();

const __prod__ = process.env.NODE_ENV === 'production';
const RedisStore = connectRedis(session);

export const sessionMiddleware =
  session({
    name: COOKIE_NAME,
    store: new RedisStore({
      client: redis as any,
      disableTouch: true,
    }),
    cookie: {
      maxAge: 1000 * 60 * 60 * 24 * 7, // 1 week
      httpOnly: true,
      sameSite: 'lax', // csrf
      secure: __prod__, // cookie only works in https,
      domain: __prod__ ? '' : undefined,
    },
    saveUninitialized: false,
    secret: process.env.SECRET as string,
    resave: true,
    rolling: true
  });
