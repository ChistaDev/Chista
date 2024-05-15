import React from 'react';
import { Box, Typography, Divider } from '@mui/material';

const TechStack = () => (
  <Box mt={3} mb={3}>
    <Divider>
      <Typography variant="h5" gutterBottom>
        TECH STACK
      </Typography>
    </Divider>

    <Box
      mt={3}
      sx={{
        display: 'flex',
        justifyContent: 'space-between',
      }}
    >
      <Box
        sx={{
          transition: 'transform 0.3s', // Add transition for smooth effect
          '&:hover': {
            transform: 'scale(1.1)', // Enlarge on hover
          },
        }}
      >
        <img
          alt="golang-badge"
          src="https://img.shields.io/badge/go-00ADD8?style=for-the-badge&logo=go&labelColor=black"
        ></img>
      </Box>
      <Box
        sx={{
          transition: 'transform 0.3s', // Add transition for smooth effect
          '&:hover': {
            transform: 'scale(1.1)', // Enlarge on hover
          },
        }}
      >
        <img
          alt="pyhton-badge"
          src="https://img.shields.io/badge/python-3776AB?style=for-the-badge&logo=python&labelColor=black"
        ></img>
      </Box>
      <Box
        sx={{
          transition: 'transform 0.3s', // Add transition for smooth effect
          '&:hover': {
            transform: 'scale(1.1)', // Enlarge on hover
          },
        }}
      >
        <img
          alt="javascript-badge"
          src="https://img.shields.io/badge/javascript-F7DF1E?style=for-the-badge&logo=javascript&labelColor=black"
        ></img>
      </Box>
      <Box
        sx={{
          transition: 'transform 0.3s', // Add transition for smooth effect
          '&:hover': {
            transform: 'scale(1.1)', // Enlarge on hover
          },
        }}
      >
        <img
          alt="react-badge"
          src="https://img.shields.io/badge/react-61DAFB?style=for-the-badge&logo=react&labelColor=black"
        ></img>
      </Box>
      <Box
        sx={{
          transition: 'transform 0.3s', // Add transition for smooth effect
          '&:hover': {
            transform: 'scale(1.1)', // Enlarge on hover
          },
        }}
      >
        <img
          alt="swagger-badge"
          src="https://img.shields.io/badge/swagger-%2385EA2D?style=for-the-badge&logo=swagger&labelColor=black"
        ></img>
      </Box>
      <Box
        sx={{
          transition: 'transform 0.3s', // Add transition for smooth effect
          '&:hover': {
            transform: 'scale(1.1)', // Enlarge on hover
          },
        }}
      >
        <img
          alt="docker-badge"
          src="https://img.shields.io/badge/docker-%232496ED?style=for-the-badge&logo=docker&labelColor=black"
        ></img>
      </Box>
    </Box>
  </Box>
);

export default TechStack;
