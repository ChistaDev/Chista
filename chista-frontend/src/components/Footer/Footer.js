import React from 'react';
import { Box, Typography, Link } from '@mui/material';
import EmailOutlinedIcon from '@mui/icons-material/EmailOutlined';
import GitHubIcon from '@mui/icons-material/GitHub';
import XIcon from '@mui/icons-material/X';
import LinkedInIcon from '@mui/icons-material/LinkedIn';

const Footer = () => (
  <Box
    sx={{
      backgroundColor: '#333',
      color: '#fff',
      py: 3,
      px: 2,
      mt: 'auto',
      textAlign: 'center',
    }}
  >
    <Typography variant="body2" gutterBottom>
      Chista Cyber Threat Intelligence Framework | © {new Date().getFullYear()}
    </Typography>
    <Typography variant="body2">
      Designed and Developed by Chista Team
    </Typography>
    <Box mt={2}>
      <Link
        href="https://twitter.com/ChistaDev"
        target="_blank"
        rel="noopener noreferrer"
        color="inherit"
        sx={{
          display: 'inline-flex',
          alignItems: 'center',
          textDecoration: 'none',
          mr: 1,
          '&:hover': {
            color: '#3373C4', // Change to bright white when hovered
          },
        }}
      >
        <XIcon sx={{ mr: 0.5 }} />
        Twitter
      </Link>
      <Link
        href="https://github.com/ChistaDev/Chista"
        target="_blank"
        rel="noopener noreferrer"
        color="inherit"
        sx={{
          display: 'inline-flex',
          alignItems: 'center',
          textDecoration: 'none',
          mr: 1,
          '&:hover': {
            color: '#3373C4', // Change to bright white when hovered
          },
        }}
      >
        <GitHubIcon sx={{ mr: 0.5 }} />
        GitHub
      </Link>
      <Link
        href="https://www.linkedin.com/company/chistadev/"
        target="_blank"
        rel="noopener noreferrer"
        color="inherit"
        sx={{
          display: 'inline-flex',
          alignItems: 'center',
          textDecoration: 'none',
          mr: 1,
          '&:hover': {
            color: '#3373C4', // Change to bright white when hovered
          },
        }}
      >
        <LinkedInIcon sx={{ mr: 0.5 }} />
        LinkedIn
      </Link>
    </Box>
    <Box mt={2}>
      <Link
        href="mailto:chistaframework@gmail.com"
        color="inherit"
        sx={{
          display: 'inline-flex',
          alignItems: 'center',
          textDecoration: 'none',
          mr: 1,
          '&:hover': {
            color: '#3373C4', // Change to bright white when hovered
          },
        }}
      >
        <EmailOutlinedIcon sx={{ mr: 0.5 }} />
        chistaframework@gmail.com
      </Link>
    </Box>
  </Box>
);

export default Footer;
