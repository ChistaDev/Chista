import React from 'react';
import { Box, Typography, Divider, Avatar } from '@mui/material';
import FF from '../../images/furkanfirat.jpg';
import FO from '../../images/furkanozturk.png';
import RB from '../../images/resulbozburun.jpeg';
import YD from '../../images/yakuphandevrez.jpg';
import YY from '../../images/yusufyildiz.jpeg';

const ProjectTeam = () => {
  const developers = [
    {
      name: 'Resul Bozburun',
      github: 'https://github.com/rbozburun',
      avatar: RB,
      role: 'AppSec Consultant',
    },
    {
      name: 'Yusuf Yıldız',
      github: 'https://github.com/grealyve',
      avatar: YY,
      role: 'Project Developer',
    },
    {
      name: 'Furkan Öztürk',
      github: 'https://github.com/furk4n0zturk',
      avatar: FO,
      role: 'Project Developer',
    },
    {
      name: 'Furkan Fırat',
      github: 'https://github.com/furkan-firat',
      avatar: FF,
      role: 'React Developer',
    },
    {
      name: 'Yakuphan Devrez',
      github: 'https://github.com/yakuphandevrez',
      avatar: YD,
      role: 'Graphic Designer',
    },
  ];
  return (
    <Box mt={5} mb={3}>
      <Divider>
        <Typography variant="h5" gutterBottom>
          PROJECT TEAM
        </Typography>
      </Divider>
      <Box mt={4} mb={4} sx={{ display: 'flex', justifyContent: 'center' }}>
        {developers.map((developer, index) => (
          <Box
            key={index}
            sx={{
              marginRight: '20px',
              display: 'flex',
              flexDirection: 'column',
              alignItems: 'center',
              transition: 'transform 0.3s', // Add transition for smooth effect
              '&:hover': {
                transform: 'scale(1.1)', // Enlarge on hover
              },
            }}
          >
            <a
              href={developer.github}
              target="_blank"
              rel="noopener noreferrer"
              style={{ textDecoration: 'none', color: 'inherit' }}
            >
              <Avatar
                alt={developer.name}
                src={developer.avatar}
                sx={{ width: 56, height: 56, marginBottom: 1 }}
              />
            </a>

            <a
              href={developer.github}
              target="_blank"
              rel="noopener noreferrer"
              style={{ textDecoration: 'none', color: 'inherit' }}
            >
              <Typography variant="body1" component="p">
                {developer.name}
              </Typography>
              <Typography variant="body2" component="p" color="textSecondary">
                {developer.role}
              </Typography>
            </a>
          </Box>
        ))}
      </Box>
    </Box>
  );
};

export default ProjectTeam;
