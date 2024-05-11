import React from 'react';
import Box from '@mui/material/Box';
import { useIsDrawerOpen } from '../contexts/IsDrawerOpenContext';
import DrawerHeader from '../components/DrawHeader/DrawHeader';
import PageHeader from '../components/PageHeader/PageHeader';
import PageDescription from '../components/PageDescription/PageDescription';
import DataLeakMonitorInputs from '../components/InputAreas/DataLeakMonitorInputs';

const DataLeakMonitor = () => {
  const { isDrawerOpen } = useIsDrawerOpen();

  return (
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
      <PageHeader pageHeader={'Data Leak Monitor'} />
      <PageDescription
        pageDescription={
          "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged."
        }
      />
      <DataLeakMonitorInputs />
    </Box>
  );
};

export default DataLeakMonitor;
