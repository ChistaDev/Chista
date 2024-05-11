import React from 'react';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import { FaTelegramPlane } from 'react-icons/fa';
import { Link } from 'react-router-dom';

const SourcesTelegramItem = () => (
  <div>
    <Link
      to="/sources/telegram-channels"
      style={{ textDecoration: 'none', color: 'inherit' }}
    >
      <ListItemButton sx={{ pl: 4 }}>
        <ListItemIcon style={{ color: '#fff', fontSize: '1.5rem' }}>
          <FaTelegramPlane />
        </ListItemIcon>
        <ListItemText
          primary="Telegram Channels"
          primaryTypographyProps={{
            noWrap: false,
            sx: { marginLeft: '-15px' },
          }}
        />
      </ListItemButton>
    </Link>
  </div>
);

export default SourcesTelegramItem;
