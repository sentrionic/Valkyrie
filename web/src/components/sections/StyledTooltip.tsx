import React from 'react';
import { Tooltip } from '@chakra-ui/react';

type Placement = 'top' | 'right';

interface StyledTooltipProps {
  label: string;
  position: Placement;
  disabled?: boolean;
}

export const StyledTooltip: React.FC<StyledTooltipProps> = ({ label, position, disabled = false, children }) => (
  <Tooltip
    hasArrow
    label={label}
    placement={position}
    isDisabled={disabled}
    bg="brandGray.darkest"
    color="white"
    fontWeight="semibold"
    py={1}
    px={3}
  >
    {children}
  </Tooltip>
);
