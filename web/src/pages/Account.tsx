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
import { Link as RLink, useHistory } from "react-router-dom";
import { AccountResponse } from "../api/response/accountresponse";
import { getAccount, logout, updateAccount } from "../api/setupAxios";
import { ColorModeSwitcher } from "../components/ColorModeSwitcher";
import { InputField } from "../components/common/InputField";
import { ChangePasswordModal } from "../components/layouts/ChangePasswordModal";
import { toErrorMap } from "../utils/toErrorMap";
import { userStore } from "../utils/userStore";
import { UserSchema } from "../utils/yup-schemas";

export const Account = () => {
  const { colorMode } = useColorMode();
  const history = useHistory();
  const toast = useToast();
  const { isOpen, onOpen, onClose } = useDisclosure();

  const { data: user } = useQuery<AccountResponse>("account", () =>
    getAccount().then((response) => response.data)
  );
  const logoutUser = userStore((state) => state.logout);
  const setUser = userStore((state) => state.setUser);

  useEffect(() => {
    if (user) setUser(user);
  }, [user, setUser]);

  const inputFile: any = useRef(null);
  const [imageUrl, setImageUrl] = useState(user?.image || "");

  const closeClicked = () => {
    history.goBack();
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
                  formData.append("image", values.image ?? imageUrl);
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
              {({ isSubmitting, setFieldValue, values }) => (
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
                        setFieldValue("image", e.currentTarget.files[0]);
                        setImageUrl(
                          URL.createObjectURL(e.currentTarget.files[0])
                        );
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
    </Flex>
  );
};
