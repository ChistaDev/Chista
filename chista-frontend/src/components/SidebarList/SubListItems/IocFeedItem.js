import React from 'react';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import { Link } from 'react-router-dom';
import VpnLockOutlinedIcon from '@mui/icons-material/VpnLockOutlined';
import { Divider } from '@mui/material';

const IocFeedItem = () => (
  <div>
    <Link
      to="/ioc/ioc-feed"
      style={{ textDecoration: 'none', color: 'inherit' }}
    >
      <ListItemButton sx={{ pl: 4 }}>
        <ListItemIcon style={{ color: '#fff' }}>
          <VpnLockOutlinedIcon />
        </ListItemIcon>
        <ListItemText
          primary="IOC Feed"
          primaryTypographyProps={{ marginLeft: '-15px' }}
        />
      </ListItemButton>
      <Divider />
    </Link>
  </div>
);

export default IocFeedItem;
