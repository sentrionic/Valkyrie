import { createParamDecorator, ExecutionContext } from '@nestjs/common';

/**
 * Returns the userId of the current user
 */
export const GetUser = createParamDecorator(
  (data: unknown, ctx: ExecutionContext) => {
    const req = ctx.switchToHttp().getRequest();
    return req.session.userId;
  },
);
