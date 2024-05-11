import React from 'react';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import HomeOutlinedIcon from '@mui/icons-material/HomeOutlined';
import { Link } from 'react-router-dom';

const HomePageItem = () => (
  <>
    <Link to="/" style={{ textDecoration: 'none', color: 'inherit' }}>
      <ListItemButton>
        <ListItemIcon style={{ color: '#fff' }}>
          <HomeOutlinedIcon />
        </ListItemIcon>
        <ListItemText
          primary="Home Page"
          primaryTypographyProps={{ marginLeft: '-7px' }}
        />
      </ListItemButton>
    </Link>
  </>
);

export default HomePageItem;
