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
import { Member } from '../../lib/api/models';
import { removeFriend } from '../../lib/api/handler/account';
import { fKey } from '../../lib/utils/querykeys';
import { useQueryClient } from 'react-query';
import { useGetFriend } from '../../lib/utils/hooks/useGetFriend';

interface IProps {
  id: string;
  isOpen: boolean;
  onClose: () => void;
}

export const RemoveFriendModal: React.FC<IProps> = ({ id, isOpen, onClose }) => {

  const cache = useQueryClient();
  const user = useGetFriend(id);

  return (
    <Modal isOpen={isOpen} onClose={onClose} isCentered>
      <ModalOverlay />

      <ModalContent bg='brandGray.light'>

        <ModalHeader textTransform={"uppercase"} fontWeight='bold' mb={0} pb={0}>
          Remove '{user?.username}'
        </ModalHeader>
        <ModalBody>
          <Text mb={"4"}>
            Are you sure you want to permanently remove <b>{user?.username}</b> from your friends?
          </Text>
        </ModalBody>

        <ModalFooter bg='brandGray.dark'>
          <Button onClick={onClose} mr={6} variant='link' fontSize={"14px"}>
            Cancel
          </Button>
          <LightMode>
            <Button
              colorScheme='red'
              fontSize={"14px"}
              onClick={async () => {
                onClose();
                const { data } = await removeFriend(id);
                if (data) {
                  cache.setQueryData<Member[]>(fKey, (d) => {
                    return d!.filter(f => f.id !== id);
                  });
                }
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
