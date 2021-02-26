import React, { useEffect, useState } from 'react';
import { Link as RLink, useHistory, useParams } from 'react-router-dom';
import { Box, Flex, Image, Link, Text } from '@chakra-ui/react';
import { joinGuild } from '../lib/api/handler/guilds';
import { Guild } from '../lib/api/models';
import { gKey } from '../lib/utils/querykeys';
import { useQueryClient } from 'react-query';

interface InviteRouter {
  link: string;
}

export const Invite: React.FC = () => {

  const { link } = useParams<InviteRouter>();
  const [errors, setErrors] = useState<string | null>(null);
  const cache = useQueryClient();
  const history = useHistory();

  useEffect(() => {
    const handleJoin = async () => {
      try {
        const { data } = await joinGuild({ link });
        if (data) {
          console.log(data);
          cache.setQueryData<Guild[]>(gKey, (old) => {
            return [...old! || [], data];
          });
          history.replace(`/channels/${data.id}/${data.default_channel_id}`);
        }
      } catch (err) {
        const status = err?.response?.status;
        if (status === 400 || status === 404) {
          setErrors(err?.response?.data?.message);
        }
        if (err?.response?.data?.errors) {
          setErrors('An error occurred. Please try again later');
        }
      }
    };
    handleJoin();
  }, [link, history, cache]);

  return (
    <Flex minHeight='100vh' align={'center'} justify={'center'} h={'full'}>
      <Box textAlign={'center'}>
        <Flex mb='4' justify='center'>
          <Image src={`${process.env.PUBLIC_URL}/logo.png`} w='80px' />
        </Flex>
        <Text>Fetching server info. Please wait.</Text>
        <Text>You will be automatically redirected.</Text>
        {errors &&
        <Box>
          <Text my={'2'} textColor={'menuRed'}>{errors}</Text>
          <Text>
            Click{' '}
            <Link
              as={RLink}
              to='/channels/me'
              color='highlight.standard'
            >
              here
            </Link>
            {' '}to return.
          </Text>
        </Box>
        }
      </Box>
    </Flex>
  );
}
