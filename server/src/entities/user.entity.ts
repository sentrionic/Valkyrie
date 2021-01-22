import { Column, Entity, JoinColumn, ManyToMany, OneToMany } from 'typeorm';
import { AbstractEntity } from './abstract.entity';
import { classToPlain, Exclude } from 'class-transformer';
import { UserResponse } from '../models/response/UserResponse';
import { Member } from './member.entity';
import { Channel } from './channel.entity';

@Entity('users')
export class User extends AbstractEntity {
  @Column('varchar', { length: 50 })
  username: string;

  @Column('varchar', { length: 255, unique: true })
  email: string;

  @Column('text')
  @Exclude()
  password: string;

  @Column('text', { nullable: true })
  image: string;

  @OneToMany(() => Member, (member) => member.user)
  guilds: Promise<Member[]>;

  @ManyToMany(() => Channel)
  @JoinColumn({ name: 'channel_member' })
  channels: Promise<Channel[]>;
  //
  // @OneToMany(() => PCMember, (pcmember) => pcmember.user)
  // pcmembers: Promise<PCMember[]>;


  toJSON(): UserResponse {
    return <UserResponse>classToPlain(this);
  }
}
