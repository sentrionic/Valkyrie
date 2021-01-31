import {
  Button,
  Flex,
  FormControl,
  FormLabel,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  Switch,
  Text,
} from '@chakra-ui/react';
import { Form, Formik } from 'formik';
import React from 'react';
import { AiOutlineLock } from 'react-icons/ai';
import { InputField } from '../common/InputField';
import { toErrorMap } from '../../lib/utils/toErrorMap';
import { createChannel } from '../../lib/api/handler/guilds';
import { ChannelSchema } from '../../lib/utils/validation/channel.schema';

interface IProps {
  guildId: string;
  isOpen: boolean;
  onClose: () => void;
}

export const CreateChannelModal: React.FC<IProps> = ({ guildId, isOpen, onClose }) => {
  return (
    <Modal isOpen={isOpen} onClose={onClose} isCentered>
      <ModalOverlay />

      <ModalContent bg='brandGray.light'>
        <Formik
          initialValues={{
            name: '',
            isPublic: true,
          }}
          validationSchema={ChannelSchema}
          onSubmit={async (values, { setErrors }) => {
            try {
              const { data } = await createChannel(guildId, values);
              if (data) {
                onClose();
              }
            } catch (err) {
              if (err?.response?.data?.errors) {
                const errors = err?.response?.data?.errors;
                setErrors(toErrorMap(errors));
              }
            }
          }
          }
        >
          {({ isSubmitting }) => (
            <Form>
              <ModalHeader textAlign='center' fontWeight='bold'>
                Create Text Channel
              </ModalHeader>
              <ModalCloseButton />
              <ModalBody pb={6}>
                <InputField label='channel name' name='name' />

                <FormControl
                  display='flex'
                  alignItems='center'
                  justifyContent='space-between'
                  mt='4'
                >
                  <FormLabel htmlFor='email-alerts' mb='0'>
                    <Flex align='center'>
                      <AiOutlineLock />
                      <Text ml='2'>Private Channel</Text>
                    </Flex>
                  </FormLabel>
                  <Switch id='isPublic' name="isPublic" />
                </FormControl>
                <Text mt='4' fontSize='14px' textColor='brandGray.accent'>
                  By making a channel private, only selected members will be
                  able to view this channel
                </Text>
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
                  isLoading={isSubmitting}
                >
                  Create Channel
                </Button>
              </ModalFooter>
            </Form>
          )}
        </Formik>
      </ModalContent>
    </Modal>
  );
};
