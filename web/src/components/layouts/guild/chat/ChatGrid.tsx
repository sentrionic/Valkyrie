import React from 'react';
import { GridItem } from '@chakra-ui/react';
import { scrollbarCss } from '../../../../lib/utils/theme';

export const ChatGrid: React.FC = ({ children }) =>
  <GridItem
    id={'chatGrid'}
    gridColumn={3}
    gridRow={'2'}
    bg='brandGray.light'
    mr='5px'
    display='flex'
    flexDirection='column-reverse'
    overflowY='auto'
    css={scrollbarCss}
  >
    {children}
  </GridItem>;
