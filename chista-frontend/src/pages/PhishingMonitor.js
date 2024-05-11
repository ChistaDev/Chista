import React from 'react';
import { Paper, Box } from '@mui/material';
import { useIsDrawerOpen } from '../contexts/IsDrawerOpenContext';
import DrawerHeader from '../components/DrawHeader/DrawHeader';
import PageHeader from '../components/PageHeader/PageHeader';
import PageDescription from '../components/PageDescription/PageDescription';
import PhishingMonitorInputs from '../components/InputAreas/PhishingMonitorInputs';

const PhishingMonitor = () => {
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
        <PageHeader pageHeader={'Phishing Monitor'} />
        <PageDescription
          pageDescription={
            'Chista detects websites created for phishing purposes and it periodically monitors them. If a new phishing website is created, you can find the result here.'
          }
        />
        <PhishingMonitorInputs />
      </Box>
    </Paper>
  );
};

export default PhishingMonitor;
