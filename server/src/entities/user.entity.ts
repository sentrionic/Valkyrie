import { Column, Entity, JoinColumn, ManyToMany, OneToMany } from 'typeorm';
import { AbstractEntity } from './abstract.entity';
import { classToPlain, Exclude, Expose } from 'class-transformer';
import { UserResponse } from '../models/response/UserResponse';
import { Member } from './member.entity';
import { Channel } from './channel.entity';
import { MemberResponse } from '../models/response/MemberResponse';

@Entity('users')
export class User extends AbstractEntity {
  @Column('varchar', { length: 50 })
  username: string;

  @Column('varchar', { length: 255, unique: true })
  @Expose({ groups: ['user'] })
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
    return <UserResponse>classToPlain(this, { groups: ['user'] });
  }

  toMember(): MemberResponse {
    return <MemberResponse>classToPlain(this);
  }
}
