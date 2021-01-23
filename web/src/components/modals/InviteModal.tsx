import {
  Button,
  FormLabel,
  Input,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  Text,
} from '@chakra-ui/react';
import React from 'react';

interface InviteModalProps {
  isOpen: boolean;
  onClose: () => void;
}

export const InviteModal: React.FC<InviteModalProps> = ({ isOpen, onClose }) => {
  return (
    <Modal isOpen={isOpen} onClose={onClose} isCentered>
      <ModalOverlay />
      <ModalContent bg='brandGray.light'>
        <ModalHeader textAlign='center' fontWeight='bold'>
          Invite Link
        </ModalHeader>
        <ModalCloseButton />
        <ModalBody>
          <FormLabel>
            <Text textTransform='uppercase'>INVITE LINK</Text>
          </FormLabel>

          <Input
            bg='brandGray.dark'
            borderColor='black'
            borderRadius='3px'
            focusBorderColor='highlight.standard'
            value='localhost:3000/asdoiasoi8dasoi'
            isReadOnly
          />
        </ModalBody>

        <ModalFooter bg='brandGray.dark'>
          <Button onClick={onClose} mr={6} variant='link'>
            Cancel
          </Button>
          <Button
            background='highlight.standard'
            color='white'
            type='submit'
            _hover={{ bg: 'highlight.hover' }}
            _active={{ bg: 'highlight.active' }}
            _focus={{ boxShadow: 'none' }}
          >
            Copy Link
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
};
