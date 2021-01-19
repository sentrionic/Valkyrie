import { BaseEntity, BeforeInsert, CreateDateColumn, Index, PrimaryColumn, UpdateDateColumn } from 'typeorm';
import { customAlphabet } from 'nanoid';

const alphabet = '0123456789';
const nanoid = customAlphabet(alphabet, 20);

export abstract class AbstractEntity extends BaseEntity {

  @PrimaryColumn()
  id: string;

  @CreateDateColumn()
  @Index()
  createdAt: Date;

  @UpdateDateColumn()
  updatedAt: Date;

  @BeforeInsert()
  async generateId() {
    this.id = await nanoid();
  }
}