import {
  Avatar,
  Box,
  Button,
  Flex,
  Input,
  LightMode,
  Modal,
  ModalBody,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  Text,
} from '@chakra-ui/react';
import React, { useState } from 'react';
import { editMessage } from '../../lib/api/handler/messages';
import { getTime } from '../../lib/utils/dateUtils';
import { Message } from '../../lib/models/message';

interface IProps {
  message: Message;
  isOpen: boolean;
  onClose: () => void;
}

export const EditMessageModal: React.FC<IProps> = ({ message, isOpen, onClose }) => {
  const [text, setNewText] = useState(message.text!);
  const [showError, toggleShow] = useState(false);

  const handleSubmit = async (): Promise<void> => {
    if (!text || !text.trim()) return;

    try {
      await editMessage(message.id, text.trim());
      onClose();
    } catch (err) {
      toggleShow(true);
    }
  };

  return (
    <Modal isOpen={isOpen} onClose={onClose} isCentered>
      <ModalOverlay />

      <ModalContent bg="brandGray.light">
        <ModalHeader fontWeight="bold" mb={0} pb={0}>
          Edit Message
        </ModalHeader>
        <ModalBody>
          <Flex alignItems="center" my="2" mr="1" justify="space-between" boxShadow="dark-lg" py={2}>
            <Flex alignItems="center">
              <Avatar h="40px" w="40px" ml="4" src={message.user.image} />
              <Box ml="3">
                <Flex alignItems="center">
                  <Text>{message.user.username}</Text>
                  <Text fontSize="12px" color="brandGray.accent" ml="3">
                    {getTime(message.createdAt)}
                  </Text>
                </Flex>
                <Input
                  value={text}
                  onChange={(e: any) => setNewText(e.target.value)}
                  bg="brandGray.dark"
                  borderColor="black"
                  borderRadius="3px"
                  focusBorderColor="none"
                />
              </Box>
            </Flex>
          </Flex>

          {showError && (
            <Text my="2" color="menuRed" align="center">
              Server Error. Try again later
            </Text>
          )}
        </ModalBody>

        <ModalFooter bg="brandGray.dark">
          <Button onClick={onClose} mr={6} variant="link" fontSize="14px" _focus={{ outline: 'none' }}>
            Cancel
          </Button>
          <LightMode>
            <Button colorScheme="green" fontSize="14px" onClick={handleSubmit}>
              Save
            </Button>
          </LightMode>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
};
