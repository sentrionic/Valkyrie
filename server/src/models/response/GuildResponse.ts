import { Channel } from '../../entities/channel.entity';

export class GuildResponse {
  id: string;
  name: string;
  channels: Channel[];
  createdAt: string;
  updatedAt: string;
}