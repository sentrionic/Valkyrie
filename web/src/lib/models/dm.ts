export interface DMChannel {
  id: string;
  user: DMMember;
}

export interface DMNotification extends DMChannel {
  count: number;
}

export interface DMMember {
  id: string;
  username: string;
  image: string;
  isOnline: boolean;
  isFriend: boolean;
}
