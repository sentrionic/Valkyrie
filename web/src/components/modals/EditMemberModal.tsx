import {
  Avatar,
  Box,
  Button,
  Divider,
  Flex,
  LightMode,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  Text,
  Tooltip,
  useDisclosure
} from '@chakra-ui/react';
import { Form, Formik } from 'formik';
import React, { useRef, useState } from 'react';
import { FaRegTrashAlt } from 'react-icons/fa';
import { InputField } from '../common/InputField';
import { toErrorMap } from '../../lib/utils/toErrorMap';
import { useGetCurrentGuild } from '../../lib/utils/hooks/useGetCurrentGuild';
import { GuildSchema } from '../../lib/utils/validation/guild.schema';
import { deleteGuild, editGuild } from '../../lib/api/handler/guilds';
import { CropImageModal } from './CropImageModal';
import { userStore } from '../../lib/stores/userStore';
import { MemberSchema } from '../../lib/utils/validation/member.schema';

interface IProps {
  isOpen: boolean;
  onClose: () => void;
}

export const EditMemberModal: React.FC<IProps> = ({ isOpen, onClose }) => {

  const current = userStore(state => state.current);

  return (
    <Modal isOpen={isOpen} onClose={onClose} isCentered>
      <ModalOverlay />
      <ModalContent bg='brandGray.light'>
        <Formik
          initialValues={{
            color: '',
            nickname: ''
          }}
          validationSchema={MemberSchema}
          onSubmit={async (values, { setErrors }) => {
            try {
              // const { data } = await editGuild(guildId, formData);
              // if (data) {
              //   resetForm();
              //   onClose();
              // }
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
              <ModalHeader fontWeight='bold' pb={0}>
                Change Appearance
              </ModalHeader>
              <ModalCloseButton />
              <ModalBody>
                <InputField label='nickname' name='nickname' />

                <Divider my={'4'} />

              </ModalBody>

              <ModalFooter bg='brandGray.dark'>
                <Button onClick={onClose} mr={6} variant='link' fontSize={'14px'}>
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
                  fontSize={'14px'}
                >
                  Save Changes
                </Button>
              </ModalFooter>
            </Form>
          )}
        </Formik>
      </ModalContent>
    </Modal>
  );
};
