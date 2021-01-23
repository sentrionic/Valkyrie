import {
  Button,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  Text,
  useToast,
} from "@chakra-ui/react";
import { Form, Formik } from "formik";
import React from "react";
import { toErrorMap } from "../../lib/utils/toErrorMap";
import { ChangePasswordSchema } from "../../lib/utils/yup-schemas";
import { InputField } from "../common/InputField";
import { changePassword } from '../../lib/api/handler/auth';

interface IProps {
  isOpen: boolean;
  onClose: () => void;
}

export const ChangePasswordModal: React.FC<IProps> = ({ isOpen, onClose }) => {
  const toast = useToast();

  return (
    <Modal isOpen={isOpen} onClose={onClose} isCentered>
      <ModalOverlay />

      <ModalContent bg="brandGray.light">
        <Formik
          initialValues={{
            currentPassword: "",
            newPassword: "",
            confirmNewPassword: "",
          }}
          validationSchema={ChangePasswordSchema}
          onSubmit={async (values, { setErrors }) => {
            try {
              const { data } = await changePassword(values);
              if (data) {
                toast({
                  title: "Changed Password",
                  description: "Successfully changed your password",
                  status: "success",
                  duration: 5000,
                  isClosable: true,
                });
                onClose();
              }
            } catch (err) {
              if (err?.response?.data?.errors) {
                const errors = err?.response?.data?.errors;
                setErrors(toErrorMap(errors));
              }
            }
          }}
        >
          {({ isSubmitting }) => (
            <Form>
              <ModalHeader textAlign="center" fontWeight="bold">
                Change your password
              </ModalHeader>
              <ModalCloseButton />
              <ModalBody pb={6}>
                <Text>Enter your current password and a new password</Text>
                <InputField
                  label="current password"
                  name="currentPassword"
                  autoComplete="password"
                  type="password"
                />

                <InputField
                  label="new password"
                  name="newPassword"
                  autoComplete="password"
                  type="password"
                />

                <InputField
                  label="confirm new password"
                  name="confirmNewPassword"
                  autoComplete="password"
                  type="password"
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
                  isLoading={isSubmitting}
                >
                  Done
                </Button>
              </ModalFooter>
            </Form>
          )}
        </Formik>
      </ModalContent>
    </Modal>
  );
};
