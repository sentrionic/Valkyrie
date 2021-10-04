import { Box, Button, Flex, Heading, Image, Link, Text } from '@chakra-ui/react';
import { Form, Formik } from 'formik';
import React from 'react';
import { Link as RLink, useHistory } from 'react-router-dom';
import { InputField } from '../components/common/InputField';
import { toErrorMap } from '../lib/utils/toErrorMap';
import { userStore } from '../lib/stores/userStore';
import { LoginSchema } from '../lib/utils/validation/auth.schema';
import { login } from '../lib/api/handler/auth';

export const Login = () => {
  const history = useHistory();
  const setUser = userStore((state) => state.setUser);

  return (
    <Flex minHeight="100vh" width="full" align="center" justifyContent="center">
      <Box px={4} width="full" maxWidth="500px" textAlign="center">
        <Flex mb="4" justify="center">
          <Image src={`${process.env.PUBLIC_URL}/logo.png`} w="80px" />
        </Flex>
        <Box p={4} borderRadius={4} background="brandGray.light">
          <Box textAlign="center">
            <Heading fontSize="24px">Welcome Back</Heading>
          </Box>
          <Box my={4} textAlign="left">
            <Formik
              initialValues={{ email: '', password: '' }}
              validationSchema={LoginSchema}
              onSubmit={async (values, { setErrors }) => {
                try {
                  const { data } = await login(values);
                  if (data) {
                    setUser(data);
                    history.push('/channels/me');
                  }
                } catch (err: any) {
                  if (err?.response?.status === 401) {
                    setErrors({ password: 'Invalid Credentials' });
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
                  <InputField label="Email" name="email" autoComplete="email" type="email" />

                  <InputField label="password" name="password" autoComplete="password" type="password" />

                  <Box mt={4}>
                    <Link as={RLink} to="/forgot-password" textColor="highlight.standard" _focus={{ outline: 'none' }}>
                      Forgot Password?
                    </Link>
                  </Box>

                  <Button
                    background="highlight.standard"
                    color="white"
                    width="full"
                    mt={4}
                    type="submit"
                    isLoading={isSubmitting}
                    _hover={{ bg: 'highlight.hover' }}
                    _active={{ bg: 'highlight.active' }}
                    _focus={{ boxShadow: 'none' }}
                  >
                    Login
                  </Button>
                  <Text mt="4">
                    Don't have an account yet?{' '}
                    <Link as={RLink} to="/register" textColor="highlight.standard" _focus={{ outline: 'none' }}>
                      Sign Up
                    </Link>
                  </Text>
                </Form>
              )}
            </Formik>
          </Box>
        </Box>
      </Box>
    </Flex>
  );
};
