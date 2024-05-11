import React from 'react';
import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import { useDarkMode } from '../../contexts/DarkModeContext';

const PageDescription = ({ pageDescription }) => {
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
        variant="body1"
        gutterBottom
        sx={{
          color: mode ? '#fff' : '#333',
          padding: '20px',
          borderRadius: '8px',
        }}
      >
        {pageDescription}
      </Typography>
    </Box>
  );
};

export default PageDescription;
