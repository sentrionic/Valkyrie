import { CanActivate, ExecutionContext, Injectable } from '@nestjs/common';
import { Member } from '../entities/member.entity';
import e from 'express';

@Injectable()
export class MemberGuard implements CanActivate {
  async canActivate(context: ExecutionContext): Promise<boolean> {
    const request: e.Request = context.switchToHttp().getRequest();
    if (!(request?.session["userId"])) return false;

    const { guildId } = request.params;

    const member = await Member.findOne({
      where: { guildId, userId: request.session["userId"] },
    });
    return !!member;
  }
}