import {
  Flex,
  GridItem,
  InputGroup,
  Input,
  Text,
} from '@chakra-ui/react';
import React, { useRef, useState } from 'react';
import { sendMessage } from '../../../lib/api/handler/messages';
import { useParams } from 'react-router-dom';
import { useQuery } from 'react-query';
import { Channel } from '../../../lib/api/models';
import { RouterProps } from '../../../routes/Routes';
import { FileUploadButton } from './FileUploadButton';
import { getSocket } from '../../../lib/api/getSocket';
import { userStore } from '../../../lib/stores/userStore';
import { channelStore } from '../../../lib/stores/channelStore';
import './css/MessageInput.css';

export const MessageInput: React.FC = () => {

  const [text, setText] = useState('');
  const [isSubmitting, setSubmitting] = useState(false);
  const [currentlyTyping, setCurrentlyTyping] = useState(false);
  const inputRef: any = useRef();

  const { guildId, channelId } = useParams<RouterProps>();
  const { data } = useQuery<Channel[]>(`channels-${guildId}`);
  const channel = data?.find(c => c.id === channelId);
  const socket = getSocket();
  const current = userStore(state => state.current);
  const isTyping = channelStore(state => state.typing);

  const handleSubmit = async () => {
    if (!text || !text.trim()) {
      return;
    }
    socket.emit('stopTyping', channelId, current?.username);
    setSubmitting(true);
    setCurrentlyTyping(false);
    const data = new FormData();
    data.append('text', text.trim());
    await sendMessage(channelId, data);
  };

  const getTypingString = (members: string[]): string => {
    switch (members.length) {
      case 1: return members[0];
      case 2: return `${members[0]} and ${members[1]}`;
      case 3: return `${members[0]}, ${members[1]} and ${members[2]}`;
      default: return "Several people";
    }
  }

  return (
    <GridItem gridColumn={3} gridRow={3} px='20px' pb={isTyping.length > 0 ? "0" : "26px"} bg='#36393f'>
      <InputGroup size='md' bg='#40444b' alignItems='center' borderRadius='8px'>
        <Input
          pl='3rem'
          name={'text'}
          placeholder={`Message #${channel?.name}`}
          border='0'
          _focus={{ border: '0' }}
          ref={inputRef}
          isDisabled={isSubmitting} value={text}
          onChange={(e) => {
            const value = e.target.value;
            if (value.trim().length === 1 && !currentlyTyping) {
              socket.emit('startTyping', channelId, current?.username);
              setCurrentlyTyping(true);
            } else if (value.length === 0) {
              socket.emit('stopTyping', channelId, current?.username);
              setCurrentlyTyping(false);
            }
            setText(value);
          }}
          onKeyDown={(e) => {
            if (e.key === 'Enter') handleSubmit().then(() => {
              setText('');
              setSubmitting(false);
              inputRef?.current?.focus();
            });
          }} />
        <FileUploadButton />
      </InputGroup>
      {isTyping.length > 0 &&
        <Flex align={'center'} fontSize={'12px'} my={1}>
          <div className="typing-indicator">
            <span />
            <span />
            <span />
          </div>
          <Text ml={'1'} fontWeight={'semibold'}>{getTypingString(isTyping)}</Text>
          <Text ml={'1'}>{isTyping.length === 1 ? 'is' : 'are'} typing... </Text>
        </Flex>
      }
    </GridItem>
  );
}
