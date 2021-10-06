import React, { useEffect, useState } from 'react';
import { Link as RLink, useHistory, useParams } from 'react-router-dom';
import { Box, Flex, Image, Link, Text } from '@chakra-ui/react';
import { joinGuild } from '../lib/api/handler/guilds';

interface InviteRouter {
  link: string;
}

export const Invite: React.FC = () => {
  const { link } = useParams<InviteRouter>();
  const [errors, setErrors] = useState<string | null>(null);
  const history = useHistory();

  useEffect(() => {
    const handleJoin = async (): Promise<void> => {
      try {
        const { data } = await joinGuild({ link });
        if (data) {
          history.replace(`/channels/${data.id}/${data.default_channel_id}`);
        }
      } catch (err: any) {
        const status = err?.response?.status;
        if (status === 400 || status === 404 || status === 500) {
          setErrors(err?.response?.data?.error?.message);
        }
      }
    };
    handleJoin();
  }, [link, history]);

  return (
    <Flex minHeight="100vh" align="center" justify="center" h="full">
      <Box textAlign="center">
        <Flex mb="4" justify="center">
          <Image src={`${process.env.PUBLIC_URL}/logo.png`} w="80px" />
        </Flex>
        <Text>Fetching server info. Please wait.</Text>
        <Text>You will be automatically redirected.</Text>
        {errors && (
          <Box>
            <Text my="2" textColor="menuRed">
              {errors}
            </Text>
            <Text>
              Click{' '}
              <Link as={RLink} to="/channels/me" color="highlight.standard" _focus={{ outline: 'none' }}>
                here
              </Link>{' '}
              to return.
            </Text>
          </Box>
        )}
      </Box>
    </Flex>
  );
};
