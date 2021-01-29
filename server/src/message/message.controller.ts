import {
  Body,
  Controller,
  Delete,
  Get,
  Param,
  Post,
  Put, Query,
  UploadedFile,
  UseGuards,
  UseInterceptors,
} from '@nestjs/common';
import { MessageService } from './message.service';
import { AuthGuard } from '../config/auth.guard';
import { GetUser } from '../config/user.decorator';
import { FileInterceptor } from '@nestjs/platform-express';
import { BufferFile } from '../types/BufferFile';
import { MessageResponse } from '../models/response/MessageResponse';
import { YupValidationPipe } from '../utils/yupValidationPipe';
import { MessageInput } from '../models/dto/MessageInput';
import { MessageSchema } from '../validation/message.schema';
import {
  ApiBody,
  ApiConsumes,
  ApiCookieAuth,
  ApiOkResponse,
  ApiOperation,
  ApiUnauthorizedResponse,
} from '@nestjs/swagger';
import { User } from '../entities/user.entity';
import { UpdateInput } from '../models/dto/UpdateInput';

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
    @Query('cursor') cursor?: string | null,
  ): Promise<MessageResponse[]> {
    return this.messageService.getMessages(channelId, userId, cursor);
  }

  @Post("/:channelId/messages")
  @UseInterceptors(FileInterceptor('file'))
  @UseGuards(AuthGuard)
  @ApiCookieAuth()
  @ApiOperation({ summary: 'Send Message' })
  @ApiOkResponse({ description: 'Message Success', type: Boolean })
  @ApiUnauthorizedResponse()
  @ApiBody({ type: MessageInput })
  @ApiConsumes('multipart/form-data')
  async createMessage(
    @GetUser() userId: string,
    @Param('channelId') channelId: string,
    @Body(new YupValidationPipe(MessageSchema)) input: MessageInput,
    @UploadedFile() file?: BufferFile,
  ): Promise<boolean> {
    return this.messageService.createMessage(userId, channelId, input, file);
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
