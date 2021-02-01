import {
  GridItem,
  InputGroup,
  Input,
  InputLeftElement,
  Icon,
} from '@chakra-ui/react';
import React, { useRef, useState } from 'react';
import { MdAddCircle } from 'react-icons/md';
import { sendMessage } from '../../../lib/api/handler/messages';
import { useParams } from 'react-router-dom';
import { useQuery } from 'react-query';
import { Channel } from '../../../lib/api/models';
import { RouterProps } from '../../../routes/Routes';

export const MessageInput: React.FC = () => {

  const [text, setText] = useState('');
  const [isSubmitting, setSubmitting] = useState(false);
  const inputRef: any = useRef();

  const { guildId, channelId } = useParams<RouterProps>();
  const { data } = useQuery<Channel[]>(`channels-${guildId}`);
  const channel = data?.find(c => c.id === channelId);

  const handleSubmit = async () => {
    if (!text || !text.trim()) {
      return;
    }
    setSubmitting(true);
    const data = new FormData();
    data.append("text", text);
    await sendMessage(channelId, data);
  };

  return (
    <GridItem gridColumn={3} gridRow={3} p='5px 20px 20px 20px' bg='#36393f'>
      <InputGroup size='md' bg='#40444b' alignItems='center' borderRadius='8px'>
        <Input
          pl='3rem'
          name={'text'}
          placeholder={`Message #${channel?.name}`}
          border='0'
          _focus={{ border: "0" }}
          ref={inputRef}
          isDisabled={isSubmitting} value={text}
          onChange={(e) => setText(e.target.value)}
          onKeyDown={(e) => {
           if (e.key === 'Enter') handleSubmit().then(() => {
             setText('');
             setSubmitting(false);
             inputRef?.current?.focus()
           });
          }} />
        <InputLeftElement _hover={{ cursor: 'pointer' }}>
          <Icon as={MdAddCircle} boxSize={'20px'} />
        </InputLeftElement>
      </InputGroup>
    </GridItem>
  );
};
