import { Column, Entity, ManyToOne, PrimaryGeneratedColumn } from 'typeorm';
import { Channel } from './channel.entity';
import { User } from './user.entity';
import { AbstractEntity } from './abstract.entity';

@Entity('messages')
export class Message extends AbstractEntity {
  @Column('text', { nullable: true })
  text: string;

  @Column('text', { nullable: true })
  url: string;

  @Column('varchar', { length: 50, nullable: true })
  filetype: string;

  @ManyToOne(() => Channel)
  channel: Channel;

  @ManyToOne(() => User, (user) => user.id)
  user: User;
}