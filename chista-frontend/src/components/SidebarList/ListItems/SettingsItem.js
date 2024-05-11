import React from 'react';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import { Link } from 'react-router-dom';
import { IoMdSettings } from 'react-icons/io';

const SettingsItem = () => (
  <>
    <Link to="/settings" style={{ textDecoration: 'none', color: 'inherit' }}>
      <ListItemButton>
        <ListItemIcon style={{ color: '#fff', fontSize: '1.5rem' }}>
          <IoMdSettings />
        </ListItemIcon>
        <ListItemText
          primary="Settings"
          primaryTypographyProps={{ marginLeft: '-7.8px' }}
        />
      </ListItemButton>
    </Link>
  </>
);

export default SettingsItem;
