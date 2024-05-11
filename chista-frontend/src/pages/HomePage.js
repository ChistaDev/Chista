import React from 'react';

import { Paper, Box, Typography } from '@mui/material';
import { useIsDrawerOpen } from '../contexts/IsDrawerOpenContext';
import DrawerHeader from '../components/DrawHeader/DrawHeader';

const HomePage = () => {
  const { isDrawerOpen } = useIsDrawerOpen();
  const contentStyle = {
    paddingLeft: isDrawerOpen ? '180px' : '0px',
    transition: 'padding-left 0.3s ease',
  };
  return (
    <Paper elevation={0} sx={{ height: '100vh' }} square>
      <Box component="main" sx={{ flexGrow: 1, p: 3, paddingLeft: 12 }}>
        <DrawerHeader />
        <Box sx={contentStyle}>
          <Typography paragraph>HOME PAGE CONTENT GOES HERE</Typography>
        </Box>
      </Box>
    </Paper>
  );
};

export default HomePage;
