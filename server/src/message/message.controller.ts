import {
  Body,
  Controller,
  Delete,
  Get,
  Param,
  Post,
  Put,
  UploadedFile,
  UseGuards,
  UseInterceptors,
} from '@nestjs/common';
import { MessageService } from './message.service';
import { AuthGuard } from '../config/auth.guard';
import { GetUser } from '../config/user.decorator';
import { FileInterceptor } from '@nestjs/platform-express';
import { BufferFile } from '../types/BufferFile';
import { Message } from '../entities/message.entity';

@Controller('channels')
export class MessageController {
  constructor(
    private readonly messageService: MessageService,
  ) {}

  // @Subscription(() => MessageSubscription, {
  //   name: MESSAGE_SUBSCRIPTION,
  //   resolve: (value) => value.channelMessage.message,
  //   filter: (payload, variables) => payload.channelId === variables.channelId,
  // })
  // @UseGuards(TeamGuard)
  // messageSubscription(@Args('channelId') channelId: string) {
  //   return this.pubSub.asyncIterator(MESSAGE_SUBSCRIPTION);
  // }
  //
  @Get("/:channelId/messages")
  @UseGuards(AuthGuard)
  async messages(
    @Param('channelId') channelId: string,
    @GetUser() userId: string,
    @Body('cursor') cursor?: string | null,
  ): Promise<Message[]> {
    return this.messageService.getMessages(cursor, channelId, userId);
  }
  //
  // @ResolveField(() => [MemberResponse!]!)
  // async user(
  //   @Parent() message: Message,
  //   @Context() ctx: MyContext,
  // ): Promise<User[]> {
  //   return await ctx.userLoader.load(message.user.id);
  // }
  //
  @Post("/:channelId/messages")
  @UseInterceptors(FileInterceptor('file'))
  @UseGuards(AuthGuard)
  async createMessage(
    @GetUser() user: string,
    @Param('channelId') channelId: string,
    @Body('text') text?: string,
    @UploadedFile() file?: BufferFile,
  ): Promise<boolean> {
    return this.messageService.createMessage(user, channelId, text, file);
  }

  @Put("/messages/:messageId")
  @UseGuards(AuthGuard)
  async editMessage(
    @GetUser() user: string,
    @Param('messageId') messageId: string,
    @Body('text') text: string,
  ): Promise<boolean> {
    return this.messageService.editMessage(user, messageId, text);
  }

  @Delete("/messages/:messageId")
  @UseGuards(AuthGuard)
  async deleteMessage(
    @GetUser() userId: string,
    @Param('messageId') messageId: string,
  ): Promise<boolean> {
    return this.messageService.deleteMessage(userId, messageId);
  }
}
