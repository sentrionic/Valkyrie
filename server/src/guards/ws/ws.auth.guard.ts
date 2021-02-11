import { CanActivate, ExecutionContext, Injectable } from '@nestjs/common';
import { Observable } from 'rxjs';

/**
 * Check if the current user is authenticated
 * using the sessionID
 */
@Injectable()
export class WsAuthGuard implements CanActivate {
  canActivate(
    context: ExecutionContext,
  ): boolean | Promise<boolean> | Observable<boolean> {
    const client = context.switchToWs().getClient();
    return !!(client?.handshake?.session["userId"]);
  }
}
