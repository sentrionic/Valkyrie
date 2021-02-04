import { Column, Entity, JoinColumn, ManyToMany, ManyToOne } from 'typeorm';
import { User } from './user.entity';
import { AbstractEntity } from './abstract.entity';
import { Guild } from './guild.entity';
import { classToPlain, Exclude } from 'class-transformer';
import { ChannelResponse } from '../models/response/ChannelResponse';

@Entity('channels')
export class Channel extends AbstractEntity {
  @Column('varchar')
  name!: string;

  @Column('boolean', { default: true })
  isPublic!: boolean;

  @Column('boolean', { default: false })
  dm!: boolean;

  @ManyToOne(() => Guild, (guild) => guild.id)
  @Exclude()
  guild!: Guild;

  @ManyToMany(() => User)
  @JoinColumn({ name: 'channel_member' })
  @Exclude()
  members!: Promise<User[]>;

  // @OneToMany(() => PCMember, (pcmember) => pcmember.channel)
  // pcmembers: Promise<PCMember[]>;

  toJson(): ChannelResponse {
    return <ChannelResponse>classToPlain(this);
  }
}
