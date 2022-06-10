import { Avatar, Flex, Icon, Text } from '@chakra-ui/react';
import React, { useCallback } from 'react';
import { MdHeadsetOff, MdMicOff } from 'react-icons/md';

interface VoiceUserVisualProps {
  username: string;
  image: string;
  isMuted: boolean;
  isDeafened: boolean;
}

const VoiceUserVisual: React.FC<VoiceUserVisualProps> = ({ username, image, isMuted, isDeafened }) => (
  <Flex
    py="1"
    px="2"
    align="center"
    w="80%"
    mr={2}
    ml={8}
    mb="1"
    justify="space-between"
    _hover={{
      bg: 'brandGray.light',
      borderRadius: '5px',
      cursor: 'pointer',
      color: '#fff',
    }}
  >
    <Flex>
      <Avatar size="xs" src={image} />
      <Text ml="2" fontSize={14}>
        {username}
      </Text>
    </Flex>
    {isDeafened && <Icon as={MdHeadsetOff} color="brandGray.accent" />}
    {isMuted && <Icon as={MdMicOff} color="brandGray.accent" />}
  </Flex>
);

interface VoiceUserProps extends VoiceUserVisualProps {
  stream?: MediaStream | null;
  muted?: boolean;
  controls?: boolean;
}

export const VoiceChannelItem: React.FC<VoiceUserProps> = ({
  username,
  image,
  stream,
  muted = false,
  controls = false,
  isMuted,
  isDeafened,
}) => {
  const refAudio = useCallback(
    (node: HTMLAudioElement) => {
      if (node && stream) {
        const audio = node;
        audio.srcObject = stream;
      }
    },
    [stream]
  );

  if (!stream) return <VoiceUserVisual image={image} username={username} isMuted={isMuted} isDeafened={isDeafened} />;

  return (
    <>
      <VoiceUserVisual image={image} username={username} isMuted={isMuted} isDeafened={isDeafened} />
      {/* eslint-disable-next-line jsx-a11y/media-has-caption */}
      <audio autoPlay ref={refAudio} muted={muted} controls={controls} />
    </>
  );
};
