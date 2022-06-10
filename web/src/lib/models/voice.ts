export interface VoiceSignal {
  userId: string;
  sdp?: RTCSessionDescription | null;
  ice?: RTCIceCandidate | null;
}

export interface VoiceResponse {
  clients: VCMember[];
  userId: string;
}

export interface VCMember {
  id: string;
  username: string;
  image: string;
  nickname?: string | null;
  isMuted: boolean;
  isDeafened: boolean;
  stream?: MediaStream | null;
}
