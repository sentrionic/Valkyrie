export interface Member {
  id: string;
  username: string;
  image: string;
  isOnline: boolean;
  isFriend: boolean;
  nickname?: string | null;
  color?: string | null;
}
