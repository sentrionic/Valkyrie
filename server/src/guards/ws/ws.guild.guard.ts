import { CanActivate, ExecutionContext, Injectable } from '@nestjs/common';
import { Member } from '../../entities/member.entity';

/**
 * Check if the current user is authenticated
 * using the sessionID and member of the guild
 */
@Injectable()
export class WsMemberGuard implements CanActivate {
  async canActivate(context: ExecutionContext): Promise<boolean> {
    const client = context.switchToWs().getClient();
    if (!client?.handshake?.session["userId"]) return false;

    const id = client?.handshake?.session["userId"];
    const guildId = context.getArgs()[1];

    if (!guildId) return false;

    const member = await Member.findOne({
      where: { guildId, userId: id },
    });
    return !!member;
  }
}
