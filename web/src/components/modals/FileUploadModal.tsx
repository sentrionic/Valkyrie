import {
  Box,
  Button,
  Modal,
  ModalBody,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
} from "@chakra-ui/react";
import React from "react";

interface IProps {
  file: File;
  isOpen: boolean;
  uploadFile: () => void;
  onClose: () => void;
}

export const FileUploadModal: React.FC<IProps> = ({
  isOpen,
  onClose,
  uploadFile,
  file
 }) => {
  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      isCentered
      closeOnOverlayClick={false}
    >
      <ModalOverlay />

      <ModalContent bg="brandGray.light">
        <ModalHeader fontWeight="bold">UPLOAD MEDIA</ModalHeader>
        <ModalBody>
          <Box h="400px" overflow="hidden" position="relative">
          </Box>
        </ModalBody>

        <ModalFooter bg="brandGray.dark">
          <Button onClick={onClose} mr={6} variant="link" fontSize={'14px'}>
            Cancel
          </Button>
          <Button
            background="highlight.standard"
            color="white"
            type="submit"
            _hover={{ bg: "highlight.hover" }}
            _active={{ bg: "highlight.active" }}
            _focus={{ boxShadow: "none" }}
            fontSize={'14px'}
          >
            Upload
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
};
