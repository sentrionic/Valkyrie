import React from 'react';
import {
  Button,
  Input,
  InputGroup,
  InputLeftAddon,
  InputRightElement,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  Text,
  useClipboard,
} from '@chakra-ui/react';
import { Form, Formik } from 'formik';
import { useQueryClient } from 'react-query';
import { userStore } from '../../lib/stores/userStore';
import { InputField } from '../common/InputField';
import { sendFriendRequest } from '../../lib/api/handler/account';
import { rKey } from '../../lib/utils/querykeys';

interface AddFriendModalProps {
  isOpen: boolean;
  onClose: () => void;
}

export const AddFriendModal: React.FC<AddFriendModalProps> = ({ isOpen, onClose }) => {
  const current = userStore((state) => state.current);
  const cache = useQueryClient();
  const { hasCopied, onCopy } = useClipboard(current?.id || '');

  return (
    <Modal isOpen={isOpen} onClose={onClose} isCentered>
      <ModalOverlay />
      <ModalContent bg="brandGray.light">
        <Formik
          initialValues={{
            id: '',
          }}
          onSubmit={async (values, { setErrors }) => {
            if (values.id === '' && values.id.length !== 20) {
              setErrors({ id: 'Enter a valid ID' });
            } else {
              try {
                const { data } = await sendFriendRequest(values.id);
                if (data) {
                  onClose();
                  await cache.invalidateQueries(rKey);
                }
              } catch (err: any) {
                if (err?.response?.data?.error) {
                  const error = err?.response?.data?.error?.message;
                  setErrors({ id: error });
                }
              }
            }
          }}
        >
          {({ isSubmitting }) => (
            <Form>
              <ModalHeader fontWeight="bold" pb="0">
                ADD FRIEND
              </ModalHeader>
              <ModalCloseButton _focus={{ outline: 'none' }} />
              <ModalBody>
                <Text mb="4">You can add a friend with their UID.</Text>
                <InputGroup mb={2}>
                  <InputLeftAddon bg="#202225" borderColor="black">
                    UID
                  </InputLeftAddon>
                  <Input
                    bg="brandGray.dark"
                    borderColor={hasCopied ? 'brandGreen' : 'black'}
                    borderRadius="3px"
                    focusBorderColor="highlight.standard"
                    value={current?.id || ''}
                    isReadOnly
                  />
                  <InputRightElement width="4.5rem">
                    <Button
                      h="1.75rem"
                      size="sm"
                      bg={hasCopied ? 'brandGreen' : 'highlight.standard'}
                      color="white"
                      _hover={{ bg: 'highlight.hover' }}
                      _active={{ bg: 'highlight.active' }}
                      _focus={{ boxShadow: 'none' }}
                      onClick={onCopy}
                    >
                      {hasCopied ? 'Copied' : 'Copy'}
                    </Button>
                  </InputRightElement>
                </InputGroup>

                <InputField label="Enter a user ID" name="id" />
              </ModalBody>
              <ModalFooter bg="brandGray.dark" mt="2">
                <Button mr={6} variant="link" onClick={onClose} fontSize="14px" _focus={{ outline: 'none' }}>
                  Cancel
                </Button>
                <Button
                  background="highlight.standard"
                  color="white"
                  type="submit"
                  _hover={{ bg: 'highlight.hover' }}
                  _active={{ bg: 'highlight.active' }}
                  _focus={{ boxShadow: 'none' }}
                  isLoading={isSubmitting}
                  fontSize="14px"
                >
                  Send Friend Request
                </Button>
              </ModalFooter>
            </Form>
          )}
        </Formik>
      </ModalContent>
    </Modal>
  );
};
