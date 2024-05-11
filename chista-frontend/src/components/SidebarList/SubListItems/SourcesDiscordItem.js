import React from 'react';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import { FaDiscord } from 'react-icons/fa6';
import { Link } from 'react-router-dom';

const SourcesDiscordItem = () => (
  <div>
    <Link
      to="/sources/discord-channels"
      style={{ textDecoration: 'none', color: 'inherit' }}
    >
      <ListItemButton sx={{ pl: 4 }}>
        <ListItemIcon style={{ color: '#fff', fontSize: '1.5rem' }}>
          <FaDiscord />
        </ListItemIcon>
        <ListItemText
          primary="Discord Channels"
          primaryTypographyProps={{
            noWrap: false,
            sx: { marginLeft: '-15px' },
          }}
        />
      </ListItemButton>
    </Link>
  </div>
);

export default SourcesDiscordItem;
