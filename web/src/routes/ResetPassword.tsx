import { Box, Button, Flex, Heading, Image, Link, Text } from '@chakra-ui/react';
import { Form, Formik } from 'formik';
import React, { useState } from 'react';
import { Link as RLink, useNavigate, useParams } from 'react-router-dom';
import { InputField } from '../components/common/InputField';
import { toErrorMap } from '../lib/utils/toErrorMap';
import { userStore } from '../lib/stores/userStore';
import { ResetPasswordSchema } from '../lib/utils/validation/auth.schema';
import { resetPassword } from '../lib/api/handler/auth';

type TokenProps = {
  token: string;
};

export const ResetPassword: React.FC = () => {
  const navigate = useNavigate();
  const { token } = useParams<keyof TokenProps>() as TokenProps;
  const [showError, setShowError] = useState(false);
  const [tokenError, setTokenError] = useState('');
  const setUser = userStore((state) => state.setUser);

  return (
    <Flex minHeight="100vh" width="full" align="center" justifyContent="center">
      <Box px={4} width="full" maxWidth="500px" textAlign="center">
        <Flex mb="4" justify="center">
          <Image src={`${process.env.PUBLIC_URL}/logo.png`} w="80px" />
        </Flex>
        <Box p={4} borderRadius={4} background="brandGray.light">
          <Box textAlign="center">
            <Heading fontSize="24px">Reset Password</Heading>
          </Box>
          <Box my={4} textAlign="left">
            <Formik
              initialValues={{
                newPassword: '',
                confirmNewPassword: '',
              }}
              validationSchema={ResetPasswordSchema}
              onSubmit={async (values, { setErrors }) => {
                try {
                  const { data } = await resetPassword({
                    ...values,
                    token,
                  });
                  if (data) {
                    setUser(data);
                    navigate('/channels/me');
                  }
                } catch (err: any) {
                  if (err?.response?.status === 500) {
                    setShowError(true);
                  } else {
                    const errors = err?.response?.data?.errors;
                    const errorMap = toErrorMap(errors);

                    if ('token' in errorMap) {
                      setTokenError(errorMap.token);
                    }
                    setErrors(errorMap);
                  }
                }
              }}
            >
              {({ isSubmitting }) => (
                <Form>
                  <InputField label="New Password" name="newPassword" autoComplete="new-password" type="password" />

                  <InputField label="Confirm New Password" name="confirmNewPassword" type="password" />

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
                    Reset Password
                  </Button>
                </Form>
              )}
            </Formik>
            {showError && (
              <Text mt="4" color="menuRed" align="center">
                Server Error. Try again later
              </Text>
            )}
            {tokenError && (
              <Flex direction="column" mt="4" justify="center" align="center">
                <Text color="menuRed">Invalid or expired token.</Text>
                <Link as={RLink} to="/forgot-password" _focus={{ outline: 'none' }}>
                  Click here to get a new token
                </Link>
              </Flex>
            )}
          </Box>
        </Box>
      </Box>
    </Flex>
  );
};
