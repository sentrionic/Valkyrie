import create from 'zustand';
import { VCMember, VoiceSignal } from '../models/voice';

type VoiceState = {
  voiceChatID: string;
  inVC: boolean;
  voiceClients: VCMember[];
  voiceJoinUserId: string;
  voiceLeaveUserId: string;
  localStream: MediaStream | null;
  connections: Map<string, RTCPeerConnection>;
  rtcSignalData: VoiceSignal;
  isMuted: boolean;
  isDeafened: boolean;
  setVoiceID: (id: string) => void;
  setVoiceClients: (clients: VCMember[]) => void;
  setInVC: (value: boolean) => void;
  setVoiceJoinUserId: (id: string) => void;
  setVoiceLeaveUserId: (id: string) => void;
  setLocalStream: (stream: MediaStream) => void;
  setRtcSignalData: (signal: VoiceSignal) => void;
  getConnection: (userId: string) => RTCPeerConnection | undefined;
  setConnection: (userId: string, connection: RTCPeerConnection) => void;
  deleteConnection: (userId: string) => void;
  clearConnections: () => void;
  setIsMuted: (value: boolean) => void;
  setIsDeafened: (value: boolean) => void;
  leaveVoice: () => void;
};

export const voiceStore = create<VoiceState>((set, get) => ({
  voiceChatID: '',
  voiceClients: [],
  inVC: false,
  voiceJoinUserId: '',
  voiceLeaveUserId: '',
  localStream: null,
  connections: new Map<string, RTCPeerConnection>(),
  rtcSignalData: { userId: '' },
  isMuted: false,
  isDeafened: false,
  setVoiceID: (id) => set({ voiceChatID: id }),
  setInVC: (value) => set({ inVC: value }),
  setVoiceClients: (clients) => set({ voiceClients: clients }),
  setVoiceJoinUserId: (id) => set({ voiceJoinUserId: id }),
  setVoiceLeaveUserId: (id) => set({ voiceLeaveUserId: id }),
  setLocalStream: (stream) => set({ localStream: stream }),
  setRtcSignalData: (signal) => set({ rtcSignalData: signal }),
  getConnection: (userId) => get().connections.get(userId),
  setConnection: (userId, connection) => get().connections.set(userId, connection),
  deleteConnection: (userId) => get().connections.delete(userId),
  clearConnections: () => {
    get().connections.forEach((v, _) => {
      v.close();
    });

    get().connections.clear();
  },
  setIsDeafened: (value) => {
    const stream = get().localStream;
    if (stream) stream.getAudioTracks()[0].enabled = !value;
    set({ isDeafened: value, localStream: stream });
  },
  setIsMuted: (value) => {
    const stream = get().localStream;
    if (stream && !get().isDeafened) stream.getAudioTracks()[0].enabled = !value;
    set({ isMuted: value, localStream: stream });
  },
  leaveVoice: () => {
    get().clearConnections();
    set({
      inVC: false,
      voiceClients: [],
      voiceJoinUserId: '',
      voiceLeaveUserId: '',
      voiceChatID: '',
    });
  },
}));
