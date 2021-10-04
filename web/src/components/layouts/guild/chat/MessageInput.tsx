import React, { useRef, useState } from 'react';
import { Flex, GridItem, InputGroup, Text, Textarea } from '@chakra-ui/react';
import ResizeTextarea from 'react-textarea-autosize';
import { useParams } from 'react-router-dom';
import { useQuery } from 'react-query';
import { FileUploadButton } from './FileUploadButton';
import { sendMessage } from '../../../../lib/api/handler/messages';
import { RouterProps } from '../../../../routes/Routes';
import { getSameSocket } from '../../../../lib/api/getSocket';
import { userStore } from '../../../../lib/stores/userStore';
import { channelStore } from '../../../../lib/stores/channelStore';
import { cKey, dmKey } from '../../../../lib/utils/querykeys';
import '../css/MessageInput.css';

export const MessageInput: React.FC = () => {
  const [text, setText] = useState('');
  const [isSubmitting, setSubmitting] = useState(false);
  const [currentlyTyping, setCurrentlyTyping] = useState(false);
  const inputRef: any = useRef();

  const { guildId, channelId } = useParams<RouterProps>();
  const qKey = guildId === undefined ? dmKey : cKey(guildId);
  const { data } = useQuery<any[]>(qKey);
  const channel = data?.find((c) => c.id === channelId);

  const socket = getSameSocket();
  const current = userStore((state) => state.current);
  const isTyping = channelStore((state) => state.typing);

  const handleSubmit = async () => {
    if (!text || !text.trim()) {
      return;
    }

    socket.send(
      JSON.stringify({
        action: 'stopTyping',
        room: channelId,
        message: current?.username,
      })
    );

    try {
      setSubmitting(true);
      setCurrentlyTyping(false);
      const data = new FormData();
      data.append('text', text.trim());
      await sendMessage(channelId, data);
    } catch (err) {}
  };

  const getTypingString = (members: string[]): string => {
    switch (members.length) {
      case 1:
        return members[0];
      case 2:
        return `${members[0]} and ${members[1]}`;
      case 3:
        return `${members[0]}, ${members[1]} and ${members[2]}`;
      default:
        return 'Several people';
    }
  };

  const getPlaceholder = (): string => {
    if (!channel) return '';

    if (channel?.user) {
      return `Message @${channel?.user.username}`;
    }
    return `Message #${channel?.name}`;
  };

  return (
    <GridItem gridColumn={3} gridRow={3} px="20px" pb={isTyping.length > 0 ? '0' : '26px'} bg="brandGray.light">
      <InputGroup size="md" bg="messageInput" alignItems="center" borderRadius="8px">
        <Textarea
          as={ResizeTextarea}
          minH="40px"
          transition="height none"
          overflow="hidden"
          w="100%"
          resize="none"
          minRows={1}
          pl="3rem"
          name={'text'}
          placeholder={getPlaceholder()}
          border="0"
          _focus={{ border: '0' }}
          ref={inputRef}
          isDisabled={isSubmitting}
          value={text}
          onChange={(e) => {
            const value = e.target.value;
            if (value.trim().length === 1 && !currentlyTyping) {
              socket.send(
                JSON.stringify({
                  action: 'startTyping',
                  room: channelId,
                  message: current?.username,
                })
              );
              setCurrentlyTyping(true);
            } else if (value.length === 0) {
              socket.send(
                JSON.stringify({
                  action: 'stopTyping',
                  room: channelId,
                  message: current?.username,
                })
              );
              setCurrentlyTyping(false);
            }
            if (value.length <= 2000) setText(value);
          }}
          onKeyDown={(e) => {
            if (e.key === 'Enter')
              handleSubmit().then(() => {
                setText('');
                setSubmitting(false);
                inputRef?.current?.focus();
              });
          }}
        />
        <FileUploadButton />
      </InputGroup>
      {isTyping.length > 0 && (
        <Flex align={'center'} fontSize={'12px'} my={1}>
          <div className="typing-indicator">
            <span />
            <span />
            <span />
          </div>
          <Text ml={'1'} fontWeight={'semibold'}>
            {getTypingString(isTyping)}
          </Text>
          <Text ml={'1'}>{isTyping.length === 1 ? 'is' : 'are'} typing... </Text>
        </Flex>
      )}
    </GridItem>
  );
};
