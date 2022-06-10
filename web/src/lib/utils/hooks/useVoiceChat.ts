/* eslint-disable react-hooks/exhaustive-deps */

import { useEffect } from 'react';
import { getSameSocket } from '../../api/getSocket';
import { VoiceSignal } from '../../models/voice';
import { userStore } from '../../stores/userStore';
import { voiceStore } from '../../stores/voiceStore';

export function useSetupVoiceChat(guildId: string): void {
  const socket = getSameSocket();

  const current = userStore((state) => state.current);

  const voiceJoinUserId = voiceStore((state) => state.voiceJoinUserId);
  const [voiceLeaveUserId, setVoiceLeaveUserId] = voiceStore((state) => [
    state.voiceLeaveUserId,
    state.setVoiceLeaveUserId,
  ]);

  const [voiceClients, setVoiceClients] = voiceStore((state) => [state.voiceClients, state.setVoiceClients]);

  const localStream = voiceStore((state) => state.localStream);
  const rtcSignalData = voiceStore((state) => state.rtcSignalData);

  const [getConnection, setConnection, deleteConnection] = voiceStore((state) => [
    state.getConnection,
    state.setConnection,
    state.deleteConnection,
  ]);

  const sendVoiceSignal = (message: VoiceSignal): void => {
    socket.send(
      JSON.stringify({
        action: 'voice-signal',
        room: guildId,
        message: { ...message },
      })
    );
  };

  // On user join, add them to the connections array and bind event handlers + create offers
  useEffect(() => {
    const onUserJoin = async (): Promise<void> => {
      // Iterate over client list
      voiceClients.forEach((user) => {
        // If the new client is not in our list
        if (!getConnection(user.id)) {
          // Add this new users Peer connection to our connections map
          setConnection(
            user.id,
            new RTCPeerConnection({
              iceServers: [{ urls: 'stun:stun.services.mozilla.com' }, { urls: 'stun:stun.l.google.com:19302' }],
            })
          );
          // Wait for peer to generate ice candidate
          getConnection(user.id)!.onicecandidate = (event: RTCPeerConnectionIceEvent) => {
            if (event.candidate !== null) {
              sendVoiceSignal({ userId: user.id, ice: event.candidate });
            }
          };

          // Event handler for peer adding their stream
          getConnection(user.id)!.ontrack = (event: RTCTrackEvent) => {
            const clients = voiceClients.map((e) => {
              if (e.id === user.id) {
                return { ...e, stream: event.streams[0] };
              }
              return e;
            });

            setVoiceClients(clients);
          };

          // Adds our local audio stream to Peer
          localStream?.getAudioTracks().forEach((track) => getConnection(user.id)!.addTrack(track, localStream!));
        }
      });

      // Create offer to new client joining if it is not the current user
      if (voiceJoinUserId !== current?.id) {
        try {
          const description = await getConnection(voiceJoinUserId)?.createOffer();

          await getConnection(voiceJoinUserId)?.setLocalDescription(description);

          sendVoiceSignal({ userId: voiceJoinUserId, sdp: getConnection(voiceJoinUserId)?.localDescription });
        } catch (err) {}
      }
    };

    if (voiceJoinUserId && voiceClients) {
      onUserJoin();
    }
  }, [voiceJoinUserId, voiceClients]);

  // New message from server, configure RTC sdp session objects
  useEffect((): void => {
    const onMessageFromServer = async (): Promise<void> => {
      const { userId, sdp, ice } = rtcSignalData;
      // Check it's not coming from the current user
      if (userId !== current?.id) {
        if (sdp) {
          try {
            await getConnection(userId)!.setRemoteDescription(new RTCSessionDescription(sdp));

            if (sdp?.type === 'offer') {
              const description = await getConnection(userId)!.createAnswer();

              // Improve audio quality
              description.sdp = description.sdp?.replace(
                'useinbandfec=1',
                'useinbandfec=1; stereo=1; maxaveragebitrate=510000'
              );

              await getConnection(userId)!.setLocalDescription(description);

              sendVoiceSignal({ userId, sdp: getConnection(userId)!.localDescription });
            }
          } catch (err) {}
        }

        if (ice) {
          try {
            await getConnection(userId)!.addIceCandidate(new RTCIceCandidate(ice));
          } catch (err) {}
        }
      }
    };

    if (voiceJoinUserId !== '') {
      onMessageFromServer();
    }
  }, [rtcSignalData, voiceJoinUserId]);

  // If a user leaves, close the peer connection and remove them from the client list
  useEffect((): void => {
    if (voiceLeaveUserId !== '') {
      // Close RTC peer connection
      getConnection(voiceLeaveUserId)?.close();
      deleteConnection(voiceLeaveUserId);
      // Remove the audio element from page
      const clients = voiceClients.filter((e) => e.id !== voiceLeaveUserId);
      setVoiceClients(clients);
      setVoiceLeaveUserId('');
    }
  }, [voiceLeaveUserId]);
}
