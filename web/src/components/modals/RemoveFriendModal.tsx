import {
  Button,
  LightMode,
  Modal,
  ModalBody,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  Text,
} from '@chakra-ui/react';
import React from 'react';
import { useQueryClient } from 'react-query';
import { removeFriend } from '../../lib/api/handler/account';
import { fKey } from '../../lib/utils/querykeys';
import { Friend } from '../../lib/models/friend';

interface IProps {
  member: Friend;
  isOpen: boolean;
  onClose: () => void;
}

export const RemoveFriendModal: React.FC<IProps> = ({ member, isOpen, onClose }) => {
  const cache = useQueryClient();

  return (
    <Modal isOpen={isOpen} onClose={onClose} isCentered>
      <ModalOverlay />

      <ModalContent bg="brandGray.light">
        <ModalHeader textTransform="uppercase" fontWeight="bold" mb={0} pb={0}>
          Remove &apos;{member?.username}&apos;
        </ModalHeader>
        <ModalBody>
          <Text mb="4">
            Are you sure you want to permanently remove <b>{member?.username}</b> from your friends?
          </Text>
        </ModalBody>

        <ModalFooter bg="brandGray.dark">
          <Button onClick={onClose} mr={6} variant="link" fontSize="14px" _focus={{ outline: 'none' }}>
            Cancel
          </Button>
          <LightMode>
            <Button
              colorScheme="red"
              fontSize="14px"
              onClick={async () => {
                onClose();
                try {
                  const { data } = await removeFriend(member.id);
                  if (data) {
                    cache.setQueryData<Friend[]>(fKey, (d) => d!.filter((f) => f.id !== member.id));
                  }
                } catch (err) {}
              }}
            >
              Remove Friend
            </Button>
          </LightMode>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
};
