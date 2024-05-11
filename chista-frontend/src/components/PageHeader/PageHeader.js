import React from 'react';
import { Box, Typography } from '@mui/material';
import { useDarkMode } from '../../contexts/DarkModeContext';

const PageHeader = ({ pageHeader }) => {
  const { mode } = useDarkMode();

  return (
    <Box
      sx={{
        width: '100%',
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
      }}
    >
      <Typography
        variant="h6"
        gutterBottom
        sx={{
          backgroundColor: mode ? 'rgba(0, 0, 0, 0.87)' : '#f0f0f0',
          color: mode ? '#fff' : '#333',
          padding: '10px',
          boxShadow: '0 2px 4px rgba(0,0,0,0.1)',
          borderRadius: '8px',
        }}
      >
        {pageHeader}
      </Typography>
    </Box>
  );
};

export default PageHeader;
