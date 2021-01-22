import { Controller, Delete, Param, UseGuards } from '@nestjs/common';
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
  // @Mutation(() => Boolean)
  // @UseGuards(AuthGuard)
  // async createMessage(
  //   @GetUser() user: User,
  //   @Args('channelId') channelId: string,
  //   @Args('text', { nullable: true }) text?: string,
  //   @Args({ name: 'file', nullable: true, type: () => GraphQLUpload })
  //     file?: FileUpload,
  // ): Promise<boolean> {
  //   return this.messageService.createMessage(user, channelId, text, file);
  // }
  //
  // @Mutation(() => DefaultResponse)
  // @UseGuards(AuthGuard)
  // async editMessage(
  //   @GetUser() user: User,
  //   @Args('messageId') messageId: string,
  //   @Args('text') text: string,
  // ): Promise<DefaultResponse> {
  //   return this.messageService.editMessage(user, messageId, text);
  // }
  //
  @Delete("/:id")
  @UseGuards(AuthGuard)
  async deleteMessage(
    @GetUser() userId: string,
    @Param('id') messageId: string,
  ): Promise<boolean> {
    return this.messageService.deleteMessage(userId, messageId);
  }
}
