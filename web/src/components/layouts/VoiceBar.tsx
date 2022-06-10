import { Box, Flex, Icon, IconButton, Text, Tooltip } from '@chakra-ui/react';
import React from 'react';
import { AiFillSignal } from 'react-icons/ai';
import { HiPhoneMissedCall } from 'react-icons/hi';
import { voiceStore } from '../../lib/stores/voiceStore';
import { useGetCurrentGuild } from '../../lib/utils/hooks/useGetCurrentGuild';

export const VoiceBar: React.FC = () => {
  const [voiceChatID, inVC, leaveVoice] = voiceStore((state) => [state.voiceChatID, state.inVC, state.leaveVoice]);
  const guild = useGetCurrentGuild(voiceChatID);

  if (!inVC) return <Box />;

  return (
    <Flex p="10px" pos="absolute" bottom="54px" w="240px" bg="accountBar" align="center" justify="space-between">
      <Box>
        <Flex>
          <Icon as={AiFillSignal} color="brandGreen" mr="1" />
          <Text color="brandGreen" fontSize={13} fontWeight="semibold">
            Voice Connected
          </Text>
        </Flex>
        <Text fontSize={12} color="brandGray.accent">
          General / {guild?.name}
        </Text>
      </Box>
      <Tooltip hasArrow label="Disconnect" placement="top" bg="brandGray.darkest" color="white">
        <IconButton
          icon={<HiPhoneMissedCall />}
          aria-label="disconnect from call"
          size="sm"
          fontSize="20px"
          variant="ghost"
          onClick={() => leaveVoice()}
        />
      </Tooltip>
    </Flex>
  );
};
