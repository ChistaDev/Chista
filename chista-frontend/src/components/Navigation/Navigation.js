import * as React from 'react';
import { styled, useTheme } from '@mui/material/styles';
import Box from '@mui/material/Box';
import MuiDrawer from '@mui/material/Drawer';
import Typography from '@mui/material/Typography';
import IconButton from '@mui/material/IconButton';
import ChevronLeftIcon from '@mui/icons-material/ChevronLeft';
import ChevronRightIcon from '@mui/icons-material/ChevronRight';
import { Link } from 'react-router-dom';
import TopNavbar from '../TopNavbar/TopNavbar';
import SidebarList from '../SidebarList/SidebarList';
import { useIsDrawerOpen } from '../../contexts/IsDrawerOpenContext';
import DrawerHeader from '../DrawHeader/DrawHeader';
import { useDarkMode } from '../../contexts/DarkModeContext';

const drawerWidth = 240;

const openedMixin = (theme) => {
  const { mode } = useDarkMode();
  return {
    width: drawerWidth,
    transition: theme.transitions.create('width', {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.enteringScreen,
    }),
    overflowX: 'hidden',
    backgroundColor: mode ? '#1E1E1E' : '#01339B',
  };
};

const closedMixin = (theme) => {
  const { mode } = useDarkMode();

  return {
    transition: theme.transitions.create('width', {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.leavingScreen,
    }),
    overflowX: 'hidden',
    width: `calc(${theme.spacing(7)} + 1px)`,
    [theme.breakpoints.up('sm')]: {
      width: `calc(${theme.spacing(8)} + 1px)`,
    },
    backgroundColor: mode ? '#1E1E1E' : '#01339B',
  };
};

const Drawer = styled(MuiDrawer, {
  shouldForwardProp: (prop) => prop !== 'open',
})(({ theme, open }) => ({
  width: drawerWidth,
  flexShrink: 0,
  whiteSpace: 'nowrap',
  boxSizing: 'border-box',
  ...(open && {
    ...openedMixin(theme),
    '& .MuiDrawer-paper': openedMixin(theme),
  }),
  ...(!open && {
    ...closedMixin(theme),
    '& .MuiDrawer-paper': closedMixin(theme),
  }),
}));

const Navigation = () => {
  const { isDrawerOpen, setIsDrawerOpen } = useIsDrawerOpen();
  const theme = useTheme();
  const toggleDrawer = () => {
    setIsDrawerOpen(!isDrawerOpen);
  };

  return (
    <Box sx={{ display: 'flex' }}>
      <TopNavbar />
      <Drawer variant="permanent" open={isDrawerOpen}>
        <DrawerHeader>
          {isDrawerOpen ? (
            <Box
              display="flex"
              justifyContent="space-between"
              alignItems={'center'}
            >
              <Link to="/" style={{ textDecoration: 'none', color: 'inherit' }}>
                <Typography
                  variant="h6"
                  noWrap
                  component="div"
                  sx={{ display: { color: '#fff', marginLeft: '10px' } }}
                >
                  CHISTA
                </Typography>
              </Link>
            </Box>
          ) : null}
          <IconButton onClick={toggleDrawer} style={{ color: '#fff' }}>
            {theme.direction === 'rtl' ? (
              <ChevronRightIcon />
            ) : (
              <ChevronLeftIcon />
            )}
          </IconButton>
        </DrawerHeader>
        <SidebarList />
      </Drawer>
    </Box>
  );
};

export default Navigation;
