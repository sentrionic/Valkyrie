import { BaseEntity, BeforeInsert, CreateDateColumn, Index, PrimaryColumn, UpdateDateColumn } from 'typeorm';
import { idGenerator } from '../utils/idGenerator';

/**
 * Represents the base entity and provides a
 * snowflake ID, createdAt and updatedAt column
 */
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
