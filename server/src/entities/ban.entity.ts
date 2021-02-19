import { Entity, JoinColumn, ManyToMany, PrimaryColumn } from 'typeorm';
import { AbstractEntity } from './abstract.entity';
import { User } from './user.entity';
import { Guild } from './guild.entity';

@Entity('bans')
export class BanEntity extends AbstractEntity {
  @PrimaryColumn()
  userId!: string;

  @PrimaryColumn()
  guildId!: string;

  @ManyToMany(() => User, (user) => user.bans, { primary: true })
  @JoinColumn({ name: 'userId' })
  user!: User;

  @ManyToMany(() => Guild, (guild) => guild.bans, { primary: true, onDelete: 'CASCADE' })
  @JoinColumn({ name: 'guildId' })
  guild!: Guild;

}
