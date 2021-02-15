import {
  Button,
  Input,
  InputGroup,
  InputRightElement,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  Text,
  useClipboard
} from '@chakra-ui/react';
import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { getInviteLink } from '../../lib/api/handler/guilds';
import { RouterProps } from '../../routes/Routes';

interface InviteModalProps {
  isOpen: boolean;
  onClose: () => void;
}

export const InviteModal: React.FC<InviteModalProps> = ({ isOpen, onClose }) => {

  const { guildId } = useParams<RouterProps>();
  const [inviteLink, setInviteLink] = useState('');
  const { hasCopied, onCopy } = useClipboard(inviteLink);

  useEffect(() => {
      if (isOpen) {
        const fetchLink = async () => {
          const { data } = await getInviteLink(guildId);
          if (data) setInviteLink(data);
        };
        fetchLink();
      }
    },
    [isOpen, setInviteLink, guildId]
  );

  return (
    <Modal isOpen={isOpen} onClose={onClose} isCentered>
      <ModalOverlay />
      <ModalContent bg='brandGray.light'>
        <ModalHeader textAlign='center' fontWeight='bold' pb={'0'}>
          Invite Link
        </ModalHeader>
        <ModalCloseButton />
        <ModalBody>

          <Text mb='4'>Share this link with others to grant access to this server</Text>
          <InputGroup>
            <Input
              bg='brandGray.dark'
              borderColor={hasCopied ? 'brandGreen' : 'black'}
              borderRadius='3px'
              focusBorderColor='highlight.standard'
              value={inviteLink}
              isReadOnly
            />
            <InputRightElement width='4.5rem'>
              <Button
                h='1.75rem' size='sm'
                bg={hasCopied ? 'brandGreen' : 'highlight.standard'}
                color='white'
                type='submit'
                _hover={{ bg: 'highlight.hover' }}
                _active={{ bg: 'highlight.active' }}
                _focus={{ boxShadow: 'none' }}
                onClick={onCopy}
              >
                {hasCopied ? 'Copied' : 'Copy'}
              </Button>
            </InputRightElement>
          </InputGroup>

          <Text my={'2'} fontSize={'12px'}>Your invite link expires in 1 day</Text>

        </ModalBody>

        <ModalFooter bg='brandGray.dark' />
      </ModalContent>
    </Modal>
  );
};
