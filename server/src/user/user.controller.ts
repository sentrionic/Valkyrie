import {
  Body,
  Controller, Delete,
  Get, Param,
  Post,
  Put,
  Req,
  Res,
  UploadedFile,
  UseGuards,
  UseInterceptors,
  ValidationPipe,
} from '@nestjs/common';
import { UserService } from './user.service';
import {
  ApiBadRequestResponse,
  ApiBody,
  ApiConsumes,
  ApiCookieAuth,
  ApiCreatedResponse,
  ApiOkResponse,
  ApiOperation, ApiTags,
  ApiUnauthorizedResponse
} from '@nestjs/swagger';
import { LoginInput } from '../models/input/LoginInput';
import e from 'express';
import { RegisterInput } from '../models/input/RegisterInput';
import { YupValidationPipe } from '../utils/yupValidationPipe';
import {
  ChangePasswordSchema,
  ForgotPasswordSchema,
  RegisterSchema,
  ResetPasswordSchema,
  UserSchema,
} from '../validation/user.schema';
import { COOKIE_NAME } from '../utils/constants';
import { ChangePasswordInput } from '../models/input/ChangePasswordInput';
import {
  ForgotPasswordInput,
  ResetPasswordInput,
} from '../models/input/ResetPasswordInput';
import { AuthGuard } from '../guards/http/auth.guard';
import { GetUser } from '../config/user.decorator';
import { UpdateInput } from '../models/input/UpdateInput';
import { FileInterceptor } from '@nestjs/platform-express';
import { BufferFile } from '../types/BufferFile';
import { UserResponse } from '../models/response/UserResponse';
import { MemberResponse } from '../models/response/MemberResponse';
import { RequestResponse } from '../models/response/RequestResponse';
import { ValidationErrors } from '../utils/serializeValidationError';

@ApiTags('Account Operation')
@Controller('account')
export class UserController {
  constructor(private userService: UserService) {}

  @Post('/register')
  @ApiOperation({ summary: 'Register Account' })
  @ApiCreatedResponse({ description: 'Newly Created User', type: UserResponse })
  @ApiBadRequestResponse({ type: ValidationErrors })
  @ApiBody({ type: RegisterInput })
  async register(
    @Body(new YupValidationPipe(RegisterSchema)) credentials: RegisterInput,
    @Req() req: e.Request,
  ): Promise<UserResponse> {
    return await this.userService.register(credentials, req);
  }

  @Post('/login')
  @ApiOperation({ summary: 'User Login' })
  @ApiOkResponse({ description: 'Current User', type: UserResponse })
  @ApiUnauthorizedResponse({ description: 'Invalid credentials' })
  @ApiBody({ type: LoginInput })
  async login(
    @Body() credentials: LoginInput,
    @Req() req: e.Request,
  ): Promise<UserResponse> {
    return await this.userService.login(credentials, req);
  }

  @Post('/logout')
  @ApiOperation({ summary: 'User Logout' })
  async logout(@Req() req: e.Request, @Res() res: e.Response): Promise<any> {
    req.session?.destroy((err) => console.log(err));
    return res.clearCookie(COOKIE_NAME).status(200).send(true);
  }

  @Put('change-password')
  @UseGuards(AuthGuard)
  @ApiCookieAuth()
  @ApiOperation({ summary: 'Change Current User Password' })
  @ApiCreatedResponse({ description: 'Successfully changed password' })
  @ApiBadRequestResponse({ type: ValidationErrors })
  @ApiBody({ type: ChangePasswordInput })
  async changePassword(
    @Body(new YupValidationPipe(ChangePasswordSchema))
    input: ChangePasswordInput,
    @GetUser() id: string,
  ): Promise<boolean> {
    return await this.userService.changePassword(input, id);
  }

  @Post('forgot-password')
  @ApiOperation({ summary: 'Forgot Password Request' })
  @ApiCreatedResponse({ description: 'Send Email' })
  @ApiBadRequestResponse({ type: ValidationErrors })
  async forgotPassword(
    @Body(new YupValidationPipe(ForgotPasswordSchema))
    { email }: ForgotPasswordInput,
  ): Promise<boolean> {
    return await this.userService.forgotPassword(email);
  }

  @Post('reset-password')
  @ApiOperation({ summary: 'Reset Password' })
  @ApiCreatedResponse({ description: 'Successfully reset password' })
  @ApiBody({ type: ResetPasswordInput })
  @ApiBadRequestResponse({ type: ValidationErrors })
  async resetPassword(
    @Body(new YupValidationPipe(ResetPasswordSchema))
    input: ResetPasswordInput,
    @Req() req: e.Request,
  ): Promise<UserResponse> {
    return await this.userService.resetPassword(input, req);
  }

  @Get()
  @UseGuards(AuthGuard)
  @ApiCookieAuth()
  @ApiOperation({ summary: 'Get Current User' })
  @ApiOkResponse({ description: 'Current user', type: UserResponse })
  @ApiUnauthorizedResponse()
  async findCurrentUser(@GetUser() id: string): Promise<UserResponse> {
    return await this.userService.findCurrentUser(id);
  }

  @Put()
  @UseGuards(AuthGuard)
  @UseInterceptors(FileInterceptor('image'))
  @ApiCookieAuth()
  @ApiOperation({ summary: 'Update Current User' })
  @ApiOkResponse({ description: 'Update Success', type: UserResponse })
  @ApiUnauthorizedResponse()
  @ApiBody({ type: UpdateInput })
  @ApiConsumes('multipart/form-data')
  async update(
    @GetUser() id: string,
    @Body(
      new YupValidationPipe(UserSchema),
      new ValidationPipe({ transform: true }),
    )
    data: UpdateInput,
    @UploadedFile() image?: BufferFile,
  ): Promise<UserResponse> {
    return await this.userService.updateUser(id, data, image);
  }

  @Get('/me/friends')
  @UseGuards(AuthGuard)
  @ApiCookieAuth()
  @ApiOperation({ summary: 'Get Current User\'s friends' })
  @ApiOkResponse({ description: 'List of users', type: [MemberResponse] })
  @ApiUnauthorizedResponse()
  async getFriends(@GetUser() id: string): Promise<MemberResponse[]> {
    return await this.userService.getFriends(id);
  }

  @Get('/me/pending')
  @UseGuards(AuthGuard)
  @ApiCookieAuth()
  @ApiOperation({ summary: 'Get Current User\'s friend requests' })
  @ApiOkResponse({ description: 'List of users', type: [RequestResponse] })
  @ApiUnauthorizedResponse()
  async getFriendRequests(@GetUser() id: string): Promise<RequestResponse[]> {
    return await this.userService.getPendingFriendRequests(id);
  }

  @Post('/:memberId/friend')
  @UseGuards(AuthGuard)
  @ApiOperation({ summary: 'Add Friend' })
  @ApiCreatedResponse({ description: 'Successfully send a friend request' })
  @ApiCookieAuth()
  async sendFriendRequest(
    @Param('memberId') memberId: string,
    @GetUser() userId: string
  ): Promise<boolean> {
    return await this.userService.sendFriendRequest(userId, memberId);
  }

  @Post('/:memberId/friend/accept')
  @UseGuards(AuthGuard)
  @ApiOperation({ summary: 'Accept Friend Request' })
  @ApiCreatedResponse({ description: 'Successfully added as friend' })
  @ApiCookieAuth()
  async addFriend(
    @Param('memberId') memberId: string,
    @GetUser() userId: string
  ): Promise<boolean> {
    return await this.userService.acceptFriendRequest(userId, memberId);
  }

  @Post('/:memberId/friend/cancel')
  @UseGuards(AuthGuard)
  @ApiOperation({ summary: 'Cancel Friend Request' })
  @ApiCreatedResponse({ description: 'Successfully canceled the request' })
  @ApiCookieAuth()
  async cancelFriendRequest(
    @Param('memberId') memberId: string,
    @GetUser() userId: string
  ): Promise<boolean> {
    return await this.userService.cancelFriendRequest(userId, memberId);
  }

  @Delete('/:memberId/friend')
  @UseGuards(AuthGuard)
  @ApiOperation({ summary: 'Remove Friend' })
  @ApiCreatedResponse({ description: 'Successfully removed friend' })
  @ApiCookieAuth()
  async removeFriend(
    @Param('memberId') memberId: string,
    @GetUser() userId: string
  ): Promise<boolean> {
    return await this.userService.removeFriend(userId, memberId);
  }
}
