export interface Guild {
  id: string;
  name: string;
  ownerId: string;
  default_channel_id: string;
  icon?: string;
  hasNotification?: boolean;
}
