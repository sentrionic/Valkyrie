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
import { AuthGuard } from '../guards/http/auth.guard';
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
import { ChannelGuard } from '../guards/http/channel.guard';

@Controller('channels')
export class MessageController {
  constructor(
    private readonly messageService: MessageService,
  ) {}

  @Get("/:channelId/messages")
  @UseGuards(ChannelGuard)
  @ApiOperation({ summary: 'Get Channel Messages' })
  @ApiUnauthorizedResponse({ description: 'Invalid credentials' })
  @ApiCookieAuth()
  @ApiOkResponse({ type: [MessageResponse] })
  async messages(
    @Param('channelId') channelId: string,
    @GetUser() userId: string,
    @Query('cursor') cursor?: string | null,
  ): Promise<MessageResponse[]> {
    return this.messageService.getMessages(channelId, userId, cursor);
  }

  @Post("/:channelId/messages")
  @UseInterceptors(FileInterceptor('file'))
  @UseGuards(ChannelGuard)
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
  @ApiCookieAuth()
  @ApiOperation({ summary: 'Edit Message' })
  @ApiOkResponse({ description: 'Edit Success', type: Boolean })
  @ApiUnauthorizedResponse()
  @ApiBody({ type: MessageInput })
  async editMessage(
    @GetUser() user: string,
    @Param('messageId') messageId: string,
    @Body(new YupValidationPipe(MessageSchema)) input: MessageInput,
  ): Promise<boolean> {
    return this.messageService.editMessage(user, messageId, input.text);
  }

  @Delete("/messages/:messageId")
  @UseGuards(AuthGuard)
  @ApiCookieAuth()
  @ApiOperation({ summary: 'Delete Message' })
  @ApiOkResponse({ description: 'Delete Success', type: Boolean })
  @ApiUnauthorizedResponse()
  async deleteMessage(
    @GetUser() userId: string,
    @Param('messageId') messageId: string,
  ): Promise<boolean> {
    return this.messageService.deleteMessage(userId, messageId);
  }
}
