import { Body, Controller, Delete, Param, Post, Put, UseGuards } from '@nestjs/common';
import { MessageService } from './message.service';
import { AuthGuard } from '../config/auth.guard';
import { GetUser } from '../config/user.decorator';

@Controller('messages')
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
  // @Query(() => [Message!]!)
  // @UseGuards(AuthGuard)
  // async messages(
  //   @Args('cursor', { nullable: true }) cursor: string,
  //   @Args('channelId') channelId: string,
  //   @GetUser() user: User,
  // ): Promise<Message[]> {
  //   return this.messageService.getMessages(cursor, channelId, user.id);
  // }
  //
  // @ResolveField(() => [MemberResponse!]!)
  // async user(
  //   @Parent() message: Message,
  //   @Context() ctx: MyContext,
  // ): Promise<User[]> {
  //   return await ctx.userLoader.load(message.user.id);
  // }
  //
  @Post("/:channelId")
  @UseGuards(AuthGuard)
  async createMessage(
    @GetUser() user: string,
    @Param('channelId') channelId: string,
    @Body('text') text?: string,
    @Args({ name: 'file', nullable: true, type: () => GraphQLUpload })
      file?: FileUpload,
  ): Promise<boolean> {
    return this.messageService.createMessage(user, channelId, text, file);
  }

  @Put("/:id")
  @UseGuards(AuthGuard)
  async editMessage(
    @GetUser() user: string,
    @Param('messageId') messageId: string,
    @Body('text') text: string,
  ): Promise<boolean> {
    return this.messageService.editMessage(user, messageId, text);
  }

  @Delete("/:id")
  @UseGuards(AuthGuard)
  async deleteMessage(
    @GetUser() userId: string,
    @Param('id') messageId: string,
  ): Promise<boolean> {
    return this.messageService.deleteMessage(userId, messageId);
  }
}
