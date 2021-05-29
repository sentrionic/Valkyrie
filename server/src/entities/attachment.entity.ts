import { Column, Entity, OneToOne } from 'typeorm';
import { AbstractEntity } from './abstract.entity';
import { classToPlain, Exclude } from 'class-transformer';
import { Message } from './message.entity';
import { AttachmentResponse } from '../models/response/AttachmentResponse';

@Entity('attachments')
export class Attachment extends AbstractEntity {
  @Column('text')
  url!: string;

  @Column('varchar', { length: 50 })
  filetype!: string;

  @Column('varchar', { nullable: true })
  filename!: string;

  @Exclude()
  @OneToOne(() => Message, (message) => message.attachment, {
    onDelete: 'CASCADE',
  })
  message: Message;

  toJSON(): AttachmentResponse {
    return <AttachmentResponse>classToPlain(this);
  }
}
