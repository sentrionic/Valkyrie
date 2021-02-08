import { BaseEntity, BeforeInsert, CreateDateColumn, Index, PrimaryColumn, UpdateDateColumn } from 'typeorm';
import { idGenerator } from '../utils/idGenerator';

export abstract class AbstractEntity extends BaseEntity {

  @PrimaryColumn()
  id!: string;

  @CreateDateColumn()
  @Index()
  createdAt!: Date;

  @UpdateDateColumn()
  updatedAt!: Date;

  @BeforeInsert()
  async generateId() {
    this.id = await idGenerator();
  }
}
