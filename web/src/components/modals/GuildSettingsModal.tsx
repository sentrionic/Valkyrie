import {
  Avatar,
  Box,
  Button,
  Divider,
  Flex,
  IconButton,
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
  useDisclosure,
} from '@chakra-ui/react';
import { Form, Formik } from 'formik';
import React, { useRef, useState } from 'react';
import { FaRegTrashAlt } from 'react-icons/fa';
import { IoPersonRemove, IoCheckmarkCircle } from 'react-icons/io5';
import { ImHammer2 } from 'react-icons/im';
import { BiUnlink } from 'react-icons/bi';
import { useQuery, useQueryClient } from 'react-query';
import { InputField } from '../common/InputField';
import { toErrorMap } from '../../lib/utils/toErrorMap';
import { useGetCurrentGuild } from '../../lib/utils/hooks/useGetCurrentGuild';
import { GuildSchema } from '../../lib/utils/validation/guild.schema';
import { deleteGuild, editGuild, getBanList, invalidateInviteLinks, unbanMember } from '../../lib/api/handler/guilds';
import { CropImageModal } from './CropImageModal';
import { Member } from '../../lib/api/models';
import { channelScrollbarCss } from '../layouts/guild/css/ChannelScrollerCSS';

interface IProps {
  guildId: string;
  isOpen: boolean;
  onClose: () => void;
}

enum SettingsScreen {
  START,
  CONFIRM,
  BANLIST,
}

export const GuildSettingsModal: React.FC<IProps> = ({ guildId, isOpen, onClose }) => {
  const guild = useGetCurrentGuild(guildId);

  const [screen, setScreen] = useState(SettingsScreen.START);
  const [isReset, setIsReset] = useState(false);
  const [showError, toggleShow] = useState(false);

  const goBack = () => setScreen(SettingsScreen.START);
  const submitClose = () => {
    setScreen(SettingsScreen.START);
    onClose();
  };

  const { isOpen: cropperIsOpen, onOpen: cropperOnOpen, onClose: cropperOnClose } = useDisclosure();

  const inputFile: any = useRef(null);
  const [imageUrl, setImageUrl] = useState<string | null>(guild?.icon || '');
  const [cropImage, setCropImage] = useState('');
  const [croppedImage, setCroppedImage] = useState<any>(null);

  const applyCrop = (file: Blob) => {
    setImageUrl(URL.createObjectURL(file));
    setCroppedImage(new File([file], 'icon', { type: 'image/jpeg' }));
    cropperOnClose();
  };

  if (!guild) return null;

  const invalidateInvites = async () => {
    try {
      const { data } = await invalidateInviteLinks(guild!.id);
      if (data) {
        setIsReset(true);
      }
    } catch (err) {}
  };

  return (
    <Modal isOpen={isOpen} onClose={onClose} isCentered>
      <ModalOverlay />
      {screen === SettingsScreen.START && (
        <ModalContent bg="brandGray.light">
          <Formik
            initialValues={{
              name: guild.name,
            }}
            validationSchema={GuildSchema}
            onSubmit={async (values, { setErrors, resetForm }) => {
              try {
                const formData = new FormData();
                formData.append('name', values.name);
                if (cropImage) {
                  formData.append('image', croppedImage);
                } else if (imageUrl) {
                  formData.append('icon', imageUrl);
                }

                const { data } = await editGuild(guildId, formData);
                if (data) {
                  resetForm();
                  onClose();
                }
              } catch (err: any) {
                if (err?.response?.status === 500) {
                  toggleShow(true);
                }
                if (err?.response?.data?.errors) {
                  const errors = err?.response?.data?.errors;
                  setErrors(toErrorMap(errors));
                }
              }
            }}
          >
            {({ isSubmitting }) => (
              <Form>
                <ModalHeader textAlign="center" fontWeight="bold" pb={0}>
                  Server Overview
                </ModalHeader>
                <ModalCloseButton _focus={{ outline: 'none' }} />
                <ModalBody>
                  <Flex mb="4" justify="center">
                    <Box textAlign={'center'}>
                      <Tooltip label="Change Icon" aria-label="Change Icon">
                        <Avatar
                          size="xl"
                          name={guild?.name[0]}
                          bg={'brandGray.darker'}
                          color={'#fff'}
                          src={imageUrl || ''}
                          _hover={{ cursor: 'pointer', opacity: 0.5 }}
                          onClick={() => inputFile.current.click()}
                        />
                      </Tooltip>
                      <Text
                        mt={'2'}
                        _hover={{
                          cursor: 'pointer',
                          color: 'brandGray.accent',
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
                      type="file"
                      name="image"
                      accept="image/*"
                      ref={inputFile}
                      hidden
                      onChange={async (e) => {
                        if (!e.currentTarget.files) return;
                        setCropImage(URL.createObjectURL(e.currentTarget.files[0]));
                        cropperOnOpen();
                      }}
                    />
                  </Flex>

                  <InputField label="server name" name="name" />

                  <Divider my={'4'} />

                  <Text fontWeight={'semibold'} mb={2}>
                    Additional Configuration
                  </Text>

                  <Flex align={'center'} justify={'space-between'} mb={'2'}>
                    <Button
                      onClick={invalidateInvites}
                      fontSize={'14px'}
                      rightIcon={isReset ? <IoCheckmarkCircle /> : <BiUnlink />}
                      colorScheme={isReset ? 'green' : 'gray'}
                    >
                      Invalidate Links
                    </Button>
                    <Button
                      onClick={() => setScreen(SettingsScreen.BANLIST)}
                      fontSize={'14px'}
                      rightIcon={<ImHammer2 />}
                    >
                      Bans
                    </Button>
                  </Flex>
                  <Flex align={'center'} justify={'space-between'} mb={'2'}>
                    <LightMode>
                      <Button
                        onClick={() => setScreen(SettingsScreen.CONFIRM)}
                        colorScheme={'red'}
                        variant="ghost"
                        fontSize={'14px'}
                        textColor={'menuRed'}
                        rightIcon={<FaRegTrashAlt />}
                      >
                        Delete Server
                      </Button>
                    </LightMode>
                  </Flex>
                  {showError && (
                    <Text my="2" color="menuRed" align="center">
                      Server Error. Try again later
                    </Text>
                  )}
                </ModalBody>

                <ModalFooter bg="brandGray.dark">
                  <Button onClick={onClose} mr={6} variant="link" fontSize={'14px'}>
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
                    fontSize={'14px'}
                  >
                    Save Changes
                  </Button>
                </ModalFooter>
              </Form>
            )}
          </Formik>
          {cropperIsOpen && (
            <CropImageModal
              isOpen={cropperIsOpen}
              onClose={cropperOnClose}
              initialImage={cropImage}
              applyCrop={applyCrop}
            />
          )}
        </ModalContent>
      )}
      {screen === SettingsScreen.CONFIRM && (
        <DeleteGuildModal goBack={goBack} submitClose={submitClose} name={guild.name} guildId={guildId} />
      )}
      {screen === SettingsScreen.BANLIST && <BanListModal goBack={goBack} guildId={guildId} />}
    </Modal>
  );
};

interface IScreenProps {
  goBack: () => void;
  submitClose: () => void;
  name: string;
  guildId: string;
}

const DeleteGuildModal: React.FC<IScreenProps> = ({ goBack, submitClose, name, guildId }) => {
  const [showError, toggleShow] = useState(false);

  const handleDelete = async () => {
    try {
      const { data } = await deleteGuild(guildId);
      if (data) {
        submitClose();
      }
    } catch (err: any) {
      if (err?.response?.status === 500) {
        toggleShow(true);
      }
    }
  };

  return (
    <ModalContent bg="brandGray.light">
      <ModalHeader fontWeight="bold" pb="0">
        Delete {name}
      </ModalHeader>
      <ModalBody pb={3}>
        <Text>
          Are you sure you want to delete <b>{name}</b>? This cannot be undone.
        </Text>

        {showError && (
          <Text my="2" color="menuRed" align="center">
            Server Error. Try again later
          </Text>
        )}
      </ModalBody>

      <ModalFooter bg="brandGray.dark">
        <Button mr={6} variant="link" onClick={goBack} fontSize={'14px'} _focus={{ outline: 'none' }}>
          Cancel
        </Button>
        <LightMode>
          <Button colorScheme="red" fontSize={'14px'} onClick={() => handleDelete()}>
            Delete Server
          </Button>
        </LightMode>
      </ModalFooter>
    </ModalContent>
  );
};

interface IBanScreenProps {
  goBack: () => void;
  guildId: string;
}

const BanListModal: React.FC<IBanScreenProps> = ({ goBack, guildId }) => {
  const key = `bans-${guildId}`;
  const { data } = useQuery(key, () => getBanList(guildId).then((response) => response.data));
  const cache = useQueryClient();

  const unbanUser = async (id: string) => {
    try {
      const { data } = await unbanMember(guildId, id);
      if (data) {
        cache.setQueryData<Member[]>(key, (d) => {
          return d!.filter((b) => b.id !== id);
        });
      }
    } catch (err) {}
  };

  return (
    <ModalContent bg="brandGray.light" maxH={'500px'}>
      <ModalHeader fontWeight="bold" pb="0">
        {data?.length} Bans
      </ModalHeader>
      <ModalBody pb={3} overflowY={'auto'} css={channelScrollbarCss}>
        <Text mb={2}>Bans are by account. Click on the icon to unban.</Text>

        {data?.map((m) => (
          <Flex
            p={'3'}
            _hover={{
              bg: 'brandGray.dark',
              borderRadius: '5px',
            }}
            align="center"
            justify="space-between"
          >
            <Flex align="center" w={'full'}>
              <Avatar size="sm" src={m.image} />
              <Text ml="2">{m.username}</Text>
            </Flex>
            <IconButton
              icon={<IoPersonRemove />}
              borderRadius="50%"
              aria-label="unban user"
              onClick={async (e) => {
                e.preventDefault();
                await unbanUser(m.id);
              }}
            />
          </Flex>
        ))}
        {data?.length === 0 && <Text>No bans yet.</Text>}
      </ModalBody>

      <ModalFooter bg="brandGray.dark">
        <Button mr={6} variant="link" onClick={goBack} fontSize={'14px'} _focus={{ outline: 'none' }}>
          Back
        </Button>
      </ModalFooter>
    </ModalContent>
  );
};
