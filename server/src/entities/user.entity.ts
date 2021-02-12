import {Column, Entity, JoinColumn, JoinTable, ManyToMany, OneToMany} from 'typeorm';
import { AbstractEntity } from './abstract.entity';
import { classToPlain, Exclude, Expose } from 'class-transformer';
import { UserResponse } from '../models/response/UserResponse';
import { Member } from './member.entity';
import { Channel } from './channel.entity';
import { MemberResponse } from '../models/response/MemberResponse';
import { PCMember } from './pcmember.entity';

@Entity('users')
export class User extends AbstractEntity {
  @Column('varchar')
  username!: string;

  @Column('varchar', { unique: true })
  @Expose({ groups: ['user'] })
  email!: string;

  @Column('text')
  @Exclude()
  password!: string;

  @Column('text', { nullable: true })
  image!: string;

  @Column({ default: true })
  isOnline!: boolean;

  @OneToMany(() => Member, (member) => member.user)
  guilds!: Member[];

  @ManyToMany(() => Channel)
  @JoinColumn({ name: 'channel_member' })
  channels!: Channel[];

  @OneToMany(() => PCMember, (pcmember) => pcmember.user)
  pcmembers: PCMember[];

  @ManyToMany(() => User, { cascade: true })
  @JoinTable({
    name: 'friends',
    joinColumn: { name: 'sender' },
    inverseJoinColumn: { name: 'receiver' }
  })
  @Exclude()
  friends!: User[];

  toJSON(): UserResponse {
    return <UserResponse>classToPlain(this, { groups: ['user'] });
  }

  toMember(userId: string = null): MemberResponse {
    const response = <MemberResponse>classToPlain(this);
    response.isFriend = (userId && this.friends?.findIndex(f => f.id === userId) !== -1);
    return response;
  }

  toFriend(): MemberResponse {
    const response = <MemberResponse>classToPlain(this);
    response.isFriend = true;
    return response;
  }
}
