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

interface IProps {
  guildId: string;
  isOpen: boolean;
  onClose: () => void;
}

enum SettingsScreen {
  START,
  CONFIRM
}

export const GuildSettingsModal: React.FC<IProps> = ({ guildId, isOpen, onClose }) => {

  const guild = useGetCurrentGuild(guildId);

  const [screen, setScreen] = useState(SettingsScreen.START);

  const goBack = () => setScreen(SettingsScreen.START);
  const submitClose = () => {
    setScreen(SettingsScreen.START);
    onClose();
  };

  const {
    isOpen: cropperIsOpen,
    onOpen: cropperOnOpen,
    onClose: cropperOnClose
  } = useDisclosure();

  const inputFile: any = useRef(null);
  const [imageUrl, setImageUrl] = useState<string | null>(guild?.icon || '');
  const [cropImage, setCropImage] = useState('');
  const [croppedImage, setCroppedImage] = useState<any>(null);

  const applyCrop = (file: Blob) => {
    setImageUrl(URL.createObjectURL(file));
    setCroppedImage(new File([file], 'icon'));
    cropperOnClose();
  };

  if (!guild) return null;

  return (
    <Modal isOpen={isOpen} onClose={onClose} isCentered>
      <ModalOverlay />
      {screen === SettingsScreen.START &&
      <ModalContent bg='brandGray.light'>
        <Formik
          initialValues={{
            name: guild.name
          }}
          validationSchema={GuildSchema}
          onSubmit={async (values, { setErrors, resetForm }) => {
            try {
              const formData = new FormData();
              formData.append("name", values.name);
              formData.append("image", croppedImage ?? imageUrl);
              const { data } = await editGuild(guildId, formData);
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
          {({ isSubmitting }) => (
            <Form>
              <ModalHeader textAlign='center' fontWeight='bold' pb={0}>
                Server Overview
              </ModalHeader>
              <ModalCloseButton />
              <ModalBody>
                <Flex mb='4' justify='center'>
                  <Box textAlign={'center'}>
                    <Tooltip label='Change Icon' aria-label='Change Icon'>
                      <Avatar
                        size='xl'
                        name={guild?.name[0]}
                        bg={'brandGray.darker'}
                        color={"#fff"}
                        src={imageUrl || ''}
                        _hover={{ cursor: 'pointer', opacity: 0.5 }}
                        onClick={() => inputFile.current.click()}
                      />
                    </Tooltip>
                    <Text
                      mt={'2'}
                      _hover={{
                        cursor: 'pointer',
                        color: 'brandGray.accent'
                      }}
                      onClick={() => {
                        setCroppedImage(null);
                        setImageUrl(null);
                      }}
                    >
                      Remove
                    </Text>
                  </Box>
                  <input
                    type='file'
                    name='image'
                    accept='image/*'
                    ref={inputFile}
                    hidden
                    onChange={async (e) => {
                      if (!e.currentTarget.files) return;
                      setCropImage(
                        URL.createObjectURL(e.currentTarget.files[0])
                      );
                      cropperOnOpen();
                    }}
                  />
                </Flex>

                <InputField label='server name' name='name' />

                <Divider my={'4'} />

                <LightMode>
                  <Button
                    onClick={() => setScreen(SettingsScreen.CONFIRM)}
                    colorScheme={'red'}
                    variant='ghost'
                    fontSize={'14px'}
                    textColor={'menuRed'}
                    rightIcon={<FaRegTrashAlt />}
                  >
                    Delete Server
                  </Button>
                </LightMode>

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
        {cropperIsOpen &&
        <CropImageModal
          isOpen={cropperIsOpen}
          onClose={cropperOnClose}
          initialImage={cropImage}
          applyCrop={applyCrop}
        />
        }
      </ModalContent>
      }
      {screen === SettingsScreen.CONFIRM &&
      <DeleteGuildModal
        goBack={goBack}
        submitClose={submitClose}
        name={guild.name}
        guildId={guildId}
      />
      }
    </Modal>
  );
};

interface IScreenProps {
  goBack: () => void;
  submitClose: () => void;
  name: string;
  guildId: string
}

const DeleteGuildModal: React.FC<IScreenProps> = ({ goBack, submitClose, name, guildId }) => {

  return (
    <ModalContent bg='brandGray.light'>
      <ModalHeader fontWeight='bold' pb='0'>
        Delete {name}
      </ModalHeader>
      <ModalBody pb={3}>
        <Text>Are you sure you want to delete <b>{name}</b>? This cannot be undone.</Text>
      </ModalBody>

      <ModalFooter bg='brandGray.dark'>
        <Button mr={6} variant='link' onClick={goBack} fontSize={'14px'}>
          Cancel
        </Button>
        <LightMode>
          <Button
            colorScheme='red'
            fontSize={'14px'}
            onClick={async () => {
              submitClose();
              await deleteGuild(guildId);
            }}
          >
            Delete Server
          </Button>
        </LightMode>
      </ModalFooter>
    </ModalContent>
  );
};
