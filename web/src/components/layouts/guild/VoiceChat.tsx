import { Box, Flex, Icon, ListItem, Text } from '@chakra-ui/react';
import React from 'react';
import { FaVolumeUp } from 'react-icons/fa';
import { useQuery } from '@tanstack/react-query';
import { useParams } from 'react-router-dom';
import { getVCMembers } from '../../../lib/api/handler/guilds';
import { useVoiceSocket } from '../../../lib/api/ws/useVoiceSocket';
import { RouterProps } from '../../../lib/models/routerProps';
import { userStore } from '../../../lib/stores/userStore';
import { voiceStore } from '../../../lib/stores/voiceStore';
import { useSetupVoiceChat } from '../../../lib/utils/hooks/useVoiceChat';
import { vcKey } from '../../../lib/utils/querykeys';
import { VoiceChannelItem } from '../../items/VoiceChannelItem';

export const VoiceChat: React.FC<{}> = () => {
  const { guildId } = useParams<keyof RouterProps>() as RouterProps;
  const current = userStore((state) => state.current);
  const key = [vcKey, guildId];

  const [voiceChatID, setVoiceID] = voiceStore((state) => [state.voiceChatID, state.setVoiceID]);
  const [inVC, setIsInVC] = voiceStore((state) => [state.inVC, state.setInVC]);
  const voiceClients = voiceStore((state) => state.voiceClients);

  const [localStream, setLocalStream] = voiceStore((state) => [state.localStream, state.setLocalStream]);

  const [isMuted, isDeafened] = voiceStore((state) => [state.isMuted, state.isDeafened]);
  const leaveVoice = voiceStore((state) => state.leaveVoice);

  const { data } = useQuery(key, () => getVCMembers(guildId).then((response) => response.data));

  useVoiceSocket();
  useSetupVoiceChat(guildId);

  const joinVoice = async (): Promise<void> => {
    if (guildId === voiceChatID) return;
    if (voiceChatID !== '') {
      leaveVoice();
    }
    await openMic();
    setIsInVC(true);
    setVoiceID(guildId);
  };

  const openMic = async (): Promise<void> => {
    try {
      // Get audio device and set better audio settings
      const result = await navigator.mediaDevices.getUserMedia({
        audio: {
          autoGainControl: false,
          channelCount: 2,
          echoCancellation: false,
          latency: 0,
          noiseSuppression: false,
          sampleRate: 48000,
          sampleSize: 16,
        },
      });
      setLocalStream(result);
    } catch (err) {}
  };

  return (
    <Box>
      <ListItem
        p="5px"
        m="0 10px"
        color="brandGray.accent"
        _hover={{
          bg: 'brandGray.light',
          borderRadius: '5px',
          cursor: 'pointer',
          color: '#fff',
        }}
        mb="2px"
        position="relative"
        onClick={() => joinVoice()}
      >
        <Flex align="center" justify="space-between">
          <Flex align="center">
            <Icon as={FaVolumeUp} color="brandGray.accent" />
            <Text ml="2">General</Text>
          </Flex>
        </Flex>
      </ListItem>
      <Box>
        {/* Current user */}
        {inVC && localStream && guildId === voiceChatID && (
          <VoiceChannelItem
            image={current!.image}
            username={current!.username}
            stream={localStream}
            controls={false}
            muted
            isMuted={isMuted}
            isDeafened={isDeafened}
          />
        )}
        {/* User is in VC, render all other clients */}
        {guildId === voiceChatID
          ? voiceClients.map(
              (e) =>
                e.id !== current?.id && (
                  <VoiceChannelItem
                    image={e.image}
                    username={e.username}
                    stream={e.stream}
                    key={e.id}
                    muted={isDeafened}
                    isMuted={e.isMuted}
                    isDeafened={e.isDeafened}
                  />
                )
            )
          : // User is not in VC, render all clients without the stream
            data?.map((e) => (
              <VoiceChannelItem
                image={e.image}
                username={e.username}
                key={e.id}
                isMuted={e.isMuted}
                isDeafened={e.isDeafened}
              />
            ))}
      </Box>
    </Box>
  );
};
