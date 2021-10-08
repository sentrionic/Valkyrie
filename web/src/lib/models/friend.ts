export interface Friend {
  id: string;
  username: string;
  image: string;
  isOnline: boolean;
}

export enum RequestType {
  OUTGOING,
  INCOMING,
}

export interface FriendRequest {
  id: string;
  username: string;
  image: string;
  type: RequestType;
}
