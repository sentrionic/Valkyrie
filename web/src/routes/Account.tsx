import {
  Avatar,
  Box,
  Button,
  Divider,
  Flex,
  Heading,
  LightMode,
  Spacer,
  Text,
  Tooltip,
  useColorMode,
  useDisclosure,
  useToast,
} from "@chakra-ui/react";
import { Form, Formik } from "formik";
import React, { useEffect, useRef, useState } from "react";
import { useQuery } from "react-query";
import { useHistory } from "react-router-dom";
import { ColorModeSwitcher } from "../components/common/ColorModeSwitcher";
import { InputField } from "../components/common/InputField";
import { ChangePasswordModal } from "../components/modals/ChangePasswordModal";
import { toErrorMap } from "../lib/utils/toErrorMap";
import { userStore } from "../lib/stores/userStore";
import { UserSchema } from "../lib/utils/validation/yup-schemas";
import { getAccount, updateAccount } from "../lib/api/handler/account";
import { AccountResponse } from "../lib/api/models";
import { logout } from "../lib/api/handler/auth";
import { CropImageModal } from "../components/modals/CropImageModal";

export const Account = () => {
  const { colorMode } = useColorMode();
  const history = useHistory();
  const toast = useToast();
  const { isOpen, onOpen, onClose } = useDisclosure();
  const {
    isOpen: cropperIsOpen,
    onOpen: cropperOnOpen,
    onClose: cropperOnClose,
  } = useDisclosure();

  const { data: user } = useQuery<AccountResponse>("account", () =>
    getAccount().then((response) => response.data)
  );
  const logoutUser = userStore((state) => state.logout);
  const setUser = userStore((state) => state.setUser);

  const inputFile: any = useRef(null);
  const [imageUrl, setImageUrl] = useState(user?.image || "");
  const [croppedImage, setCroppedImage] = useState<any>(null);

  const closeClicked = () => {
    history.goBack();
  };

  const applyCrop = (file: Blob) => {
    setImageUrl(URL.createObjectURL(file));
    setCroppedImage(new File([file], "avatar"));
    cropperOnClose();
  };

  const logoutClicked = async () => {
    const { data } = await logout();
    if (data) {
      logoutUser();
      history.replace("/");
    }
  };

  if (!user) return null;

  return (
    <Flex minHeight="100vh" width="full" align="center" justifyContent="center">
      <Box px={4} width="full" maxWidth="500px">
        <Flex mb="4" justify="center">
          <Heading fontSize="24px">MY ACCOUNT</Heading>
        </Flex>
        <Box p={4} borderRadius={4} background="brandGray.light">
          <Box>
            <Formik
              initialValues={{
                email: user.email,
                username: user.username,
                image: null,
              }}
              validationSchema={UserSchema}
              onSubmit={async (values, { setErrors }) => {
                try {
                  const formData = new FormData();
                  formData.append("email", values.email);
                  formData.append("username", values.username);
                  formData.append("image", croppedImage ?? imageUrl);
                  const { data } = await updateAccount(formData);
                  if (data) {
                    setUser(data);
                    toast({
                      title: "Account Updated.",
                      description: "Successfully updated your account",
                      status: "success",
                      duration: 5000,
                      isClosable: true,
                    });
                  }
                } catch (err) {
                  if (err?.response?.data?.errors) {
                    const errors = err?.response?.data?.errors;
                    setErrors(toErrorMap(errors));
                  }
                }
              }}
            >
              {({ isSubmitting, values }) => (
                <Form>
                  <Flex mb="4" justify="center">
                    <Tooltip label="Change Avatar" aria-label="Change Avatar">
                      <Avatar
                        size="xl"
                        name={user?.username}
                        src={imageUrl || user?.image}
                        _hover={{ cursor: "pointer", opacity: 0.5 }}
                        onClick={() => inputFile.current.click()}
                      />
                    </Tooltip>
                    <input
                      type="file"
                      name="image"
                      accept="image/*"
                      ref={inputFile}
                      hidden
                      onChange={async (e) => {
                        if (!e.currentTarget.files) return;
                        setImageUrl(
                          URL.createObjectURL(e.currentTarget.files[0])
                        );
                        cropperOnOpen();
                      }}
                    />
                  </Flex>
                  <Box my={4}>
                    <InputField
                      value={values.email}
                      type="email"
                      placeholder="Email"
                      label="Email"
                      name="email"
                      autoComplete="email"
                    />

                    <InputField
                      value={values.username}
                      placeholder="Username"
                      label="Username"
                      name="username"
                      autoComplete="username"
                    />

                    <Flex mt={8}>
                      <Flex align="center">
                        <ColorModeSwitcher />
                        <Text ml="2">
                          Use {colorMode === "light" ? "Dark" : "Light"} Mode
                        </Text>
                      </Flex>
                      <Spacer />
                      <Button
                        mr={4}
                        colorScheme="white"
                        variant="outline"
                        onClick={closeClicked}
                      >
                        Close
                      </Button>

                      <LightMode>
                        <Button
                          type="submit"
                          colorScheme="green"
                          isLoading={isSubmitting}
                        >
                          Update
                        </Button>
                      </LightMode>
                    </Flex>
                  </Box>
                </Form>
              )}
            </Formik>
          </Box>
          <Divider my={"4"} />
          <Flex>
            <Heading fontSize="18px">PASSWORD AND AUTHENTICATION</Heading>
          </Flex>
          <Flex mt="4">
            <Button
              background="highlight.standard"
              color="white"
              type="submit"
              _hover={{ bg: "highlight.hover" }}
              _active={{ bg: "highlight.active" }}
              _focus={{ boxShadow: "none" }}
              onClick={onOpen}
            >
              Change Password
            </Button>

            <Spacer />
            <Button colorScheme="red" variant="outline" onClick={logoutClicked}>
              Logout
            </Button>
          </Flex>
        </Box>
      </Box>
      <ChangePasswordModal isOpen={isOpen} onClose={onClose} />
      <CropImageModal
        isOpen={cropperIsOpen}
        onClose={cropperOnClose}
        initialImage={imageUrl}
        applyCrop={applyCrop}
      />
    </Flex>
  );
};
