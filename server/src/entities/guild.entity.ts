import { Column, Entity, ManyToOne, OneToMany } from 'typeorm';
import { AbstractEntity } from './abstract.entity';
import { Member } from './member.entity';
import { classToPlain, Exclude } from 'class-transformer';
import { GuildResponse } from '../models/response/GuildResponse';
import { BanEntity } from './ban.entity';

@Entity('guilds')
export class Guild extends AbstractEntity {
  @Column('varchar')
  name!: string;

  @Column('varchar')
  ownerId!: string;

  @ManyToOne(() => Member, (member) => member.guild)
  @Exclude()
  members!: Member[];

  @OneToMany(() => BanEntity, (bans) => bans.guild)
  @Exclude()
  bans!: BanEntity[];

  @Column('varchar', { nullable: true })
  icon?: string;

  @Column("simple-array", { default: [] })
  inviteLinks: string[];

  toJson(): GuildResponse {
    return <GuildResponse>classToPlain(this);
  }
}
