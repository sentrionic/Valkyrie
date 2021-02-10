import {
  Avatar,
  Box,
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
  Text
} from '@chakra-ui/react';
import { Form, Formik } from 'formik';
import React, { useState } from 'react';
import { useQuery } from 'react-query';
import { AiOutlineLock } from 'react-icons/ai';
import { InputField } from '../common/InputField';
import { toErrorMap } from '../../lib/utils/toErrorMap';
import { getGuildMembers } from '../../lib/api/handler/guilds';
import { ChannelSchema } from '../../lib/utils/validation/channel.schema';
import { CUIAutoComplete } from 'chakra-ui-autocomplete';
import { mKey } from '../../lib/utils/querykeys';
import { createChannel } from '../../lib/api/handler/channel';

interface IProps {
  guildId: string;
  isOpen: boolean;
  onClose: () => void;
}

interface Item {
  value: string;
  label: string;
  image: string;
}

export const CreateChannelModal: React.FC<IProps> = ({ guildId, isOpen, onClose }) => {

  const key = mKey(guildId);
  const { data } = useQuery(key, () =>
    getGuildMembers(guildId).then(response => response.data)
  );

  const members: Item[] = [];
  const [selectedItems, setSelectedItems] = useState<Item[]>([]);

  data?.map(m => members.push({ label: m.username, value: m.id, image: m.image }));

  const handleCreateItem = (item: Item) => {
    setSelectedItems((curr) => [...curr, item]);
  };

  const handleSelectedItemsChange = (selectedItems?: Item[]) => {
    if (selectedItems) {
      setSelectedItems(selectedItems);
    }
  };

  const ListItem = (selected: Item) => {
    return (
      <Flex align="center">
        <Avatar mr={2} size="sm" src={selected.image} />
        <Text textColor={'#000'}>{selected.label}</Text>
      </Flex>
    )
  }

  return (
    <Modal isOpen={isOpen} onClose={onClose} isCentered>
      <ModalOverlay />
      <ModalContent bg='brandGray.light'>
        <Formik
          initialValues={{
            name: '',
            isPublic: true
          }}
          validationSchema={ChannelSchema}
          onSubmit={async (values, { setErrors, resetForm }) => {
            try {
              const ids: string[] = [];
              selectedItems.map(i => ids.push(i.value));
              const { data } = await createChannel(guildId, { ...values, members: ids });
              if (data) {
                resetForm();
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
          {({ isSubmitting, setFieldValue, values }) => (
            <Form>
              <ModalHeader textAlign='center' fontWeight='bold'>
                Create Text Channel
              </ModalHeader>
              <ModalCloseButton />
              <ModalBody>
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
                  <Switch onChange={(e) => {
                    setFieldValue('isPublic', !e.target.checked);
                  }} />
                </FormControl>
                <Text mt='4' fontSize='14px' textColor='brandGray.accent'>
                  By making a channel private, only selected members will be
                  able to view this channel
                </Text>
                {!values.isPublic &&
                  <Box mt={"2"} pb={0}>
                  <CUIAutoComplete
                    label="Who can access this channel"
                    placeholder=""
                    onCreateItem={handleCreateItem}
                    items={members}
                    selectedItems={selectedItems}
                    itemRenderer={ListItem}
                    onSelectedItemsChange={(changes) =>
                      handleSelectedItemsChange(changes.selectedItems)
                    }
                  />
                  </Box>
                }
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
