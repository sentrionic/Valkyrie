import {
  Icon,
  InputLeftElement,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalHeader,
  ModalOverlay,
  Progress,
  Text,
  useDisclosure,
} from '@chakra-ui/react';
import React, { useRef, useState } from 'react';
import { MdAddCircle } from 'react-icons/md';
import { useParams } from 'react-router-dom';
import { sendMessage } from '../../../../lib/api/handler/messages';
import { FileSchema } from '../../../../lib/utils/validation/message.schema';
import { StyledTooltip } from '../../../sections/StyledTooltip';
import { RouterProps } from '../../../../lib/models/routerProps';

export const FileUploadButton: React.FC = () => {
  const { channelId } = useParams<keyof RouterProps>() as RouterProps;
  const { isOpen, onOpen, onClose } = useDisclosure();

  const inputFile: any = useRef(null);
  const [isSubmitting, setSubmitting] = useState(false);
  const [progress, setProgress] = useState(0);
  const [errors, setErrors] = useState({});
  const disable = process.env.NODE_ENV === 'production';

  const closeModal = (): void => {
    setErrors({});
    setProgress(0);
    onClose();
  };

  const handleSubmit = async (file: File): Promise<void> => {
    if (!file) return;
    setSubmitting(true);

    try {
      await FileSchema.validate({ file });
    } catch (err: any) {
      setErrors(err.errors);
      onOpen();
      return;
    }

    const data = new FormData();
    data.append('file', file);
    await sendMessage(channelId, data, (event: any) => {
      const loaded = Math.round((100 * event.loaded) / event.total);
      setProgress(loaded);
      if (loaded >= 100) setProgress(0);
    });
  };

  return (
    <StyledTooltip disabled={!disable} label="File Upload is disabled on the test site" position="top">
      <InputLeftElement
        color="iconColor"
        _hover={{
          cursor: 'pointer',
          color: '#fcfcfc',
        }}
        onClick={() => inputFile.current.click()}
      >
        <Icon as={MdAddCircle} boxSize="20px" />
        <input
          type="file"
          ref={inputFile}
          hidden
          disabled={isSubmitting || disable}
          onChange={async (e) => {
            if (!e.currentTarget.files) return;
            handleSubmit(e.currentTarget.files[0]).then(() => {
              setSubmitting(false);
              e.target.value = '';
            });
          }}
        />
        {errors && (
          <Modal size="sm" isOpen={isOpen} onClose={closeModal} isCentered>
            <ModalOverlay />
            <ModalContent bg="brandGray.light" textAlign="center">
              <ModalHeader pb="0">Error Uploading File</ModalHeader>
              <ModalCloseButton _focus={{ outline: 'none' }} />
              <ModalBody>
                <Text mb="2">
                  Reason: <>{errors}</>
                </Text>
                <Text>Max file size is 5.00 MB</Text>
                <Text>Only Images and mp3 allowed</Text>
              </ModalBody>
            </ModalContent>
          </Modal>
        )}
        {progress > 0 && (
          <Modal size="sm" isOpen={progress > 0} closeOnOverlayClick={false} onClose={closeModal} isCentered>
            <ModalContent bg="brandGray.darker" textAlign="center">
              <ModalHeader pb="0">Upload Progress</ModalHeader>
              <ModalBody>
                <Progress hasStripe isAnimated value={progress} />
              </ModalBody>
            </ModalContent>
          </Modal>
        )}
      </InputLeftElement>
    </StyledTooltip>
  );
};
