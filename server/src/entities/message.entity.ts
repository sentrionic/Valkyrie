import { Column, Entity, JoinColumn, ManyToOne, OneToOne } from 'typeorm';
import { AbstractEntity } from './abstract.entity';
import { Channel } from './channel.entity';
import { User } from './user.entity';
import { Attachment } from './attachment.entity';
import { classToPlain, Exclude } from 'class-transformer';
import { MessageResponse } from '../models/response/MessageResponse';

@Entity('messages')
export class Message extends AbstractEntity {
  @Column('text', { nullable: true })
  text!: string;

  @ManyToOne(() => Channel, { onDelete: 'CASCADE' })
  @Exclude()
  channel!: Channel;

  @ManyToOne(() => User, (user) => user.id)
  @Exclude()
  user!: User;

  @OneToOne(
    () => Attachment,
    attachment => attachment.message,
    { nullable: true }
  )
  @JoinColumn()
  attachment?: Attachment;

  toJSON(userId: string): MessageResponse {
    const response = <MessageResponse>classToPlain(this);
    response.user = this.user.toMember(userId);
    return response;
  }
}
