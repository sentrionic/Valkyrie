import { Box } from '@chakra-ui/react';
import React from 'react';

export const NotificationIndicator: React.FC = () => (
  <Box w="8px" h="8px" bg="white" position="absolute" borderRadius="0 4px 4px 0" ml="-4px" mt="20px" left={0} />
);

export const ChannelNotificationIndicator: React.FC = () => (
  <Box w="8px" h="8px" bg="white" position="absolute" borderRadius="0 4px 4px 0" ml="-4px" mt="8px" left="-10px" />
);

export const ActiveGuildPill: React.FC = () => (
  <Box w="8px" h="40px" bg="white" position="absolute" borderRadius="0 4px 4px 0" ml="-4px" left={0} mt="4px" />
);

export const HoverGuildPill: React.FC = () => (
  <Box w="8px" h="24px" bg="white" position="absolute" borderRadius="0 4px 4px 0" ml="-4px" left={0} mt="12px" />
);
