import { extendTheme } from '@chakra-ui/react';
import { mode } from '@chakra-ui/theme-tools';

const config: any = {
  initialColorMode: 'dark',
};

const styles = {
  global: (props: any) => ({
    body: {
      bg: mode('gray.100', '#1b1c1d')(props),
    },
  }),
};

const colors = {
  highlight: {
    standard: '#7289da',
    hover: '#677bc4',
    active: '#5b6eae',
  },
  brandGray: {
    accent: '#8e9297',
    active: '#393c43',
    light: '#36393f',
    dark: '#303339',
    darker: '#202225',
    darkest: '#18191c',
    hover: '#32353b',
  },
  brandGreen: '#43b581',
  labelGray: '#72767d',
  menuRed: '#f04747',
  brandBorder: '#1A202C',
  accountBar: '#292b2f',
  memberList: '#2f3136',
  iconColor: '#b9bbbe',
  messageInput: '#40444b',
};

const fonts = {
  body: "'Open Sans', sans-serif",
};

const customTheme = extendTheme({
  colors,
  config,
  styles,
  fonts,
});

export default customTheme;

export const scrollbarCss = {
  '&::-webkit-scrollbar': {
    width: '8px',
  },
  '&::-webkit-scrollbar-track': {
    background: '#2f3136',
    width: '10px',
  },
  '&::-webkit-scrollbar-thumb': {
    background: 'brandGray.darker',
    borderRadius: '18px',
  },
};
