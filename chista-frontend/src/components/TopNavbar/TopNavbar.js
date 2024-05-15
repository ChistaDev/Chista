import React from 'react';
import { styled } from '@mui/material/styles';
import HomeOutlinedIcon from '@mui/icons-material/HomeOutlined';
import { AppBar, Box, Toolbar, IconButton } from '@mui/material';
import MenuIcon from '@mui/icons-material/Menu';
import { Link } from 'react-router-dom';
import { useIsDrawerOpen } from '../../contexts/IsDrawerOpenContext';
import Logo from '../../images/chista-banner-145x40.png';
import { useDarkMode } from '../../contexts/DarkModeContext';
import DarkModeButton from '../DarkModeButton/DarkModeButton';

const drawerWidth = 240;

const StyledAppBar = styled(AppBar, {
  shouldForwardProp: (prop) => prop !== 'open' && prop !== 'mode',
})(({ theme, open, mode }) => ({
  zIndex: theme.zIndex.drawer + 1,
  transition: theme.transitions.create(['width', 'margin'], {
    easing: theme.transitions.easing.sharp,
    duration: theme.transitions.duration.leavingScreen,
  }),
  ...(open && {
    marginLeft: drawerWidth,
    width: `calc(100% - ${drawerWidth}px)`,
    transition: theme.transitions.create(['width', 'margin'], {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.enteringScreen,
    }),
  }),
  backgroundColor: mode ? '#121212' : '#3373C4',
}));

const TopNavbar = () => {
  const { isDrawerOpen, setIsDrawerOpen } = useIsDrawerOpen();

  const toggleDrawer = () => {
    setIsDrawerOpen(!isDrawerOpen);
  };

  const { mode } = useDarkMode();

  return (
    <Box sx={{ flexGrow: 1 }}>
      <StyledAppBar position="fixed" open={isDrawerOpen} mode={mode}>
        <Toolbar>
          <IconButton
            color="inherit"
            aria-label="open drawer"
            onClick={toggleDrawer}
            edge="start"
            sx={{
              marginRight: 5,
              ...(isDrawerOpen && { display: 'none' }),
            }}
          >
            <MenuIcon />
          </IconButton>
          {!isDrawerOpen ? (
            <>
              <Link
                to="/"
                sx={{
                  textDecoration: 'none',
                  color: 'inherit',
                  alignItems: 'center',
                }}
              >
                <img src={Logo} alt="Logo" style={{ height: '40px' }} />
              </Link>
            </>
          ) : null}
          <Box sx={{ flexGrow: 1 }} />

          <Box sx={{ display: { xs: 'none', md: 'flex' } }}>
            <IconButton
              size="large"
              aria-label="show 4 new mails"
              color="inherit"
              sx={{
                backgroundColor: 'transparent',
                color: 'inherit',
                '&:hover': {
                  backgroundColor: 'transparent',
                },
              }}
            >
              <Link to="/" style={{ textDecoration: 'none', color: 'inherit' }}>
                <HomeOutlinedIcon />
              </Link>
            </IconButton>

            <DarkModeButton />
          </Box>
        </Toolbar>
      </StyledAppBar>
    </Box>
  );
};

export default TopNavbar;
