import React, { useState } from 'react';
import { Paper, Box, Typography, Button } from '@mui/material';
import { useIsDrawerOpen } from '../contexts/IsDrawerOpenContext';
import { useDarkMode } from '../contexts/DarkModeContext';
import DrawerHeader from '../components/DrawHeader/DrawHeader';
import Banner from '../images/chista-banner-580-160.png';
import ModulesModal from '../components/ModulesModal/ModulesModal';
import TechStack from '../components/TechStack/TechStack';
import ProjectTeam from '../components/ProjectTeam/ProjectTeam';
import Footer from '../components/Footer/Footer';

const HomePage = () => {
  const { isDrawerOpen } = useIsDrawerOpen();
  const contentStyle = {
    transition: 'padding-left 0.3s ease',
  };

  const { mode } = useDarkMode();

  const [isModalOpen, setIsModalOpen] = useState(false);

  const handleDocsButtonClick = () => {
    window.open(
      'https://github.com/ChistaDev/Chista/blob/main/README.md',
      '_blank'
    );
  };

  const openModal = () => {
    setIsModalOpen(true);
  };

  const closeModal = () => {
    setIsModalOpen(false);
  };

  return (
    <>
      <Paper elevation={0} sx={{ height: 'auto' }} square>
        <Box
          sx={{ paddingLeft: isDrawerOpen ? '240px' : '0px', ...contentStyle }}
        >
          <DrawerHeader />
          <Box maxWidth="960px" mx="auto" p={3} textAlign="center">
            <Box>
              <img
                name="version"
                src="https://img.shields.io/badge/version-2.0.0-blue"
              ></img>
            </Box>
            <img
              src={Banner}
              alt="CHISTA Banner"
              style={{
                paddingTop: '25px',
                paddingBottom: '25px',
                maxWidth: '580px',
                maxHeight: '160px',
                margin: '0 auto',
              }}
            />
            <Typography variant="h4" component="h1" pb={3}>
              Chista | Open Source Threat Intelligence Framework
            </Typography>
            <Box
              mt={3}
              mb={1}
              sx={{
                display: 'flex',
                justifyContent: 'center',
              }}
            >
              <Button
                variant="contained"
                color="primary"
                onClick={handleDocsButtonClick}
                sx={{ marginRight: '50px' }}
              >
                Read the Docs
              </Button>
              <Button
                variant="contained"
                onClick={() => openModal()}
                sx={{
                  backgroundColor: mode ? 'rgb(192, 50, 50)' : 'rgb(162,0,0)',
                  color: mode ? 'black' : 'white',
                  '&:hover': {
                    backgroundColor: 'rgb(142, 0, 0) ',
                  },
                }}
              >
                Discover Modules
              </Button>
            </Box>

            <Box
              mt={4}
              mb={4}
              sx={{ display: 'flex', justifyContent: 'center' }}
            >
              <ul
                style={{
                  display: 'flex',
                  listStyle: 'none',
                  justifyContent: 'center',
                }}
              >
                <li>
                  <iframe
                    className="github-btn"
                    src="//ghbtns.com/github-btn.html?user=ChistaDev&amp;repo=Chista&amp;type=watch&amp;count=true"
                    frameBorder="0"
                    scrolling="0"
                    width="100px"
                    height="20px"
                  ></iframe>
                </li>
                <li id="version"></li>
                <li>
                  <iframe
                    className="github-btn"
                    src="//ghbtns.com/github-btn.html?user=ChistaDev&amp;repo=Chista&amp;type=fork&amp;count=true"
                    frameBorder="0"
                    scrolling="0"
                    width="100px"
                    height="20px"
                  ></iframe>
                </li>
              </ul>
            </Box>

            <Typography variant="body1" component="p" pb={3}>
              Chista is an Open Source Cyber Threat Intelligence (CTI) Framework
              designed to help users understand, predict and defend against
              cyber threats.
            </Typography>

            {/* YT VIDEO */}
            <Box
              mt={3}
              mb={4}
              sx={{ display: 'flex', justifyContent: 'center' }}
            >
              <iframe
                width="560"
                height="315"
                src="https://www.youtube.com/"
                frameBorder="0"
                allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
                allowFullScreen
              ></iframe>
            </Box>

            <Typography variant="body1" component="p" pb={3}>
              It helps its users understand cyber threats by using data
              collected from various sources. This data includes IOCs, data
              leaks, phishing campaigns, threat group activities and CTI
              sources. By analyzing this data, Chista helps users understand the
              existence, prevalence, trends and probability of cyber threats.
            </Typography>

            <TechStack />

            <ProjectTeam />

            <ModulesModal isModalOpen={isModalOpen} closeModal={closeModal} />
          </Box>
        </Box>
        <Footer />
      </Paper>
    </>
  );
};

export default HomePage;
