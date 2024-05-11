import React from 'react';
import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import CircularProgress from '@mui/material/CircularProgress';

const Loading = ({ isButtonClicked }) => (
  <Box display="flex" justifyContent="center">
    <Box variant="contained" color="primary" disabled={isButtonClicked}>
      {isButtonClicked ? (
        <>
          <Box
            display="flex"
            flexDirection="column"
            alignItems="center"
            gap={5}
            sx={{ maxWidth: '1316px', marginTop: '60px', marginRight: '160px' }}
          >
            <CircularProgress size={36} style={{}} />
            <Typography variant="body1">
              Searching for phishing URLs Please do not leave this page...
            </Typography>
          </Box>
        </>
      ) : (
        <></>
      )}
    </Box>
  </Box>
);

export default Loading;
