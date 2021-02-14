import React from 'react';
import { Tooltip } from '@chakra-ui/react';

type Placement = "top" | "right";

interface StyledTooltipProps {
  label: string;
  position: Placement;
}

export const StyledTooltip: React.FC<StyledTooltipProps> = ({ label, position , children}) =>
  <Tooltip
    hasArrow
    label={label}
    placement={position}
    bg={'#18191c'}
    color={"white"}
    fontWeight={"semibold"}
    py={1}
    px={3}
  >
    {children}
  </Tooltip>
