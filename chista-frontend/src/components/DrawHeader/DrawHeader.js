import React from 'react';
import { styled } from '@mui/material/styles';

const StyledDrawerHeader = styled('div')(({ theme }) => ({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  padding: theme.spacing(0, 1),
  // necessary for content to be below app bar
  ...theme.mixins.toolbar,
}));

const DrawerHeader = ({ children }) => (
  <StyledDrawerHeader>{children}</StyledDrawerHeader>
);

export default DrawerHeader;
