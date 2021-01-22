import {
  Avatar,
  Box,
  Button,
  Flex,
  FormLabel,
  GridItem,
  Heading,
  Icon,
  IconButton,
  Input,
  ListItem,
  Menu,
  MenuButton,
  MenuItem,
  MenuList,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  Text,
  UnorderedList,
  useDisclosure,
} from "@chakra-ui/react";
import React, { useRef } from "react";
import { FaHashtag } from "react-icons/fa";
import { RiSettings5Fill } from "react-icons/ri";
import { FiChevronDown, FiX } from "react-icons/fi";
import { InputField } from "../../common/InputField";
import { CreateChannelModal } from "../../modals/CreateChannelModal";

export const Channels: React.FC = () => {
  const { isOpen, onOpen, onClose } = useDisclosure();
  const {
    isOpen: channelIsOpen,
    onOpen: channelOpen,
    onClose: channelClose,
  } = useDisclosure();

  return (
    <GridItem gridColumn={2} gridRow={"1 / 4"} bg="brandGray.dark">
      <Menu placement="bottom-end">
        {({ isOpen }) => (
          <>
            <Flex
              justify="space-between"
              align="center"
              boxShadow="md"
              p="10px"
            >
              <Heading fontSize="20px">Harmony</Heading>
              <MenuButton>
                <Icon as={!isOpen ? FiChevronDown : FiX} />
              </MenuButton>
            </Flex>
            <MenuList bg="#18191c">
              <MenuItem onClick={channelOpen}>Create Channel</MenuItem>
              <MenuItem onClick={onOpen}>Invite People</MenuItem>
              <MenuItem>Leave Server</MenuItem>
            </MenuList>
          </>
        )}
      </Menu>
      <Modal isOpen={isOpen} onClose={onClose} isCentered>
        <ModalOverlay />
        <ModalContent bg="brandGray.light">
          <ModalHeader textAlign="center" fontWeight="bold">
            Invite Link
          </ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            <FormLabel>
              <Text textTransform="uppercase">INVITE LINK</Text>
            </FormLabel>

            <Input
              bg="brandGray.dark"
              borderColor="black"
              borderRadius="3px"
              focusBorderColor="highlight.standard"
              value="localhost:3000/asdoiasoi8dasoi"
              isReadOnly
            />
          </ModalBody>

          <ModalFooter bg="brandGray.dark">
            <Button onClick={onClose} mr={6} variant="link">
              Cancel
            </Button>
            <Button
              background="highlight.standard"
              color="white"
              type="submit"
              _hover={{ bg: "highlight.hover" }}
              _active={{ bg: "highlight.active" }}
              _focus={{ boxShadow: "none" }}
            >
              Copy Link
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
      <CreateChannelModal onClose={channelClose} isOpen={channelIsOpen} />
      <UnorderedList listStyleType="none" ml="0" mt="4">
        <ListItem
          p="5px"
          m="0 10px"
          _hover={{ bg: "#36393f", borderRadius: "5px", cursor: "pointer" }}
        >
          <Flex align="center">
            <FaHashtag />
            <Text ml="2">general</Text>
          </Flex>
        </ListItem>
      </UnorderedList>
      <Flex
        p="10px"
        pos="absolute"
        bottom="0"
        w="240px"
        bg="#292b2f"
        align="center"
        justify="space-between"
      >
        <Flex align="center">
          <Avatar size="sm" />
          <Text ml="2">Username</Text>
        </Flex>
        <IconButton
          icon={<RiSettings5Fill />}
          aria-label="settings"
          size="sm"
          fontSize="20px"
          variant="ghost"
        />
      </Flex>
    </GridItem>
  );
};
