import React from 'react';
import { Box, Flex, Link, Stack, Text } from '@chakra-ui/react';
import { AiOutlineApi, AiOutlineGithub } from 'react-icons/ai';
import { SiSocketdotio } from 'react-icons/si';
import { IconType } from 'react-icons';
import { StyledTooltip } from './StyledTooltip';

type FooterLinkProps = {
  icon?: IconType;
  href?: string;
  label?: string;
};

const FooterLink: React.FC<FooterLinkProps> = ({ icon, href, label }) => (
  <StyledTooltip label={label!} position="top">
    <Link display="inline-block" href={href} aria-label={label} isExternal mx={2}>
      <Box as={icon} width="24px" height="24px" color="gray.400" _hover={{ color: 'gray.100' }} />
    </Link>
  </StyledTooltip>
);

const links = [
  {
    icon: AiOutlineGithub,
    label: 'GitHub',
    href: 'https://github.com/sentrionic/Valkyrie',
  },
  {
    icon: AiOutlineApi,
    label: 'REST API',
    href: `${process.env.REACT_APP_API!}/swagger/index.html`,
  },
  {
    icon: SiSocketdotio,
    label: 'Websocket',
    href: `${process.env.REACT_APP_API!}`,
  },
];

export const Footer: React.FC = () => (
  <Flex bottom={0} as="footer" align="center" justify="center" w="100%" p={8}>
    <Box textAlign="center">
      <Text fontSize="xl">
        <span>Valkyrie | 2022</span>
      </Text>
      <Text>This is not a real commercial app.</Text>
      <Stack mt={2} isInline justify="center" spacing="6">
        {links.map((link) => (
          // eslint-disable-next-line react/jsx-props-no-spreading
          <FooterLink key={link.href} {...link} />
        ))}
      </Stack>
    </Box>
  </Flex>
);
