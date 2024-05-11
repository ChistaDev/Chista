import React from 'react';
import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import { useIsDrawerOpen } from '../contexts/IsDrawerOpenContext';
import DrawerHeader from '../components/DrawHeader/DrawHeader';

const NotFoundPage = () => {
  const { isDrawerOpen } = useIsDrawerOpen();
  const contentStyle = {
    paddingLeft: isDrawerOpen ? '180px' : '0px',
    transition: 'padding-left 0.3s ease',
  };
  return (
    <Box component="main" sx={{ flexGrow: 1, p: 3, paddingLeft: 12 }}>
      <DrawerHeader />
      <Box sx={contentStyle}>
        <Typography paragraph>404 PAGE NOT FOUND</Typography>
      </Box>
    </Box>
  );
};

export default NotFoundPage;
