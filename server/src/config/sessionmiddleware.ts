import * as session from 'express-session';
import { COOKIE_NAME, PRODUCTION } from '../utils/constants';
import { redis } from './redis';
import * as connectRedis from 'connect-redis';
import { config } from 'dotenv';
config();

const RedisStore = connectRedis(session);

export const sessionMiddleware = session({
  name: COOKIE_NAME,
  store: new RedisStore({
    client: redis as any,
    disableTouch: true,
  }),
  cookie: {
    maxAge: 1000 * 60 * 60 * 24 * 7, // 1 week
    httpOnly: true,
    sameSite: 'lax', // csrf
    secure: PRODUCTION, // cookie only works in https,
    domain: PRODUCTION ? '.valkyrieapp.xyz' : undefined,
  },
  saveUninitialized: false,
  secret: process.env.SECRET as string,
  resave: true,
  rolling: true,
});
