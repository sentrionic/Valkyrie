import { Column, CreateDateColumn, Entity, JoinTable, ManyToMany, ManyToOne, OneToMany } from 'typeorm';
import { User } from './user.entity';
import { AbstractEntity } from './abstract.entity';
import { Guild } from './guild.entity';
import { PCMember } from './pcmember.entity';

@Entity('channels')
export class Channel extends AbstractEntity {
  @Column('varchar')
  name!: string;

  @Column('boolean', { default: true })
  isPublic!: boolean;

  @Column('boolean', { default: false })
  dm!: boolean;

  @ManyToOne(
    () => Guild,
    (guild) => guild.id,
    {
      nullable: true,
      onDelete: 'CASCADE'
    }
  )
  guild!: Guild;

  @ManyToMany(() => User, { onDelete: 'CASCADE' })
  @JoinTable({
    name: 'channel_member',
    joinColumn: {
      name: 'channels',
      referencedColumnName: 'id'
    },
    inverseJoinColumn: {
      name: 'users',
      referencedColumnName: 'id'
    }
  })
  members!: User[];

  @OneToMany(
    () => PCMember,
    (pcmember) => pcmember.channel
  )
  pcmembers: PCMember[];

  @CreateDateColumn()
  lastActivity?: string;
}
