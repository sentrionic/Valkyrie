import { Column, Entity, ManyToOne } from 'typeorm';
import { AbstractEntity } from './abstract.entity';
import { Member } from './member.entity';
import { classToPlain, Exclude } from 'class-transformer';
import { GuildResponse } from '../models/response/GuildResponse';

@Entity('guilds')
export class Guild extends AbstractEntity {
  @Column('varchar', { length: 50 })
  name: string;

  @Column('varchar')
  ownerId: string;

  @ManyToOne(() => Member, (member) => member.guild)
  @Exclude()
  members: Promise<Member[]>;

  toJson(): GuildResponse {
    return <GuildResponse>classToPlain(this);
  }
}