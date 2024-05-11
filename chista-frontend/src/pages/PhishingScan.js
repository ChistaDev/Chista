import React from 'react';
import { Paper, Box } from '@mui/material';
import { useIsDrawerOpen } from '../contexts/IsDrawerOpenContext';
import DrawerHeader from '../components/DrawHeader/DrawHeader';
import PageHeader from '../components/PageHeader/PageHeader';
import PageDescription from '../components/PageDescription/PageDescription';
import PhishingScanInputs from '../components/InputAreas/PhishingScanInputs';

const PhishingScan = () => {
  const { isDrawerOpen } = useIsDrawerOpen();

  return (
    <Paper elevation={0} sx={{ height: '100vh' }} square>
      <Box
        component="main"
        sx={{
          flexGrow: 1,
          p: 3,
          paddingLeft: isDrawerOpen ? '280px' : '100px',
          transition: 'padding-left 0.3s ease',
        }}
      >
        <DrawerHeader />
        <PageHeader pageHeader={'Phishing Scan'} />
        <PageDescription
          pageDescription={
            'Chista detects websites created for phishing purposes and provides users with a feed in this direction.'
          }
        />
        <PhishingScanInputs />
      </Box>
    </Paper>
  );
};

export default PhishingScan;
