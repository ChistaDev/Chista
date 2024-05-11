import React from 'react';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import ForumOutlinedIcon from '@mui/icons-material/ForumOutlined';
import { Link } from 'react-router-dom';

const SourcesForumsItem = () => (
  <div>
    <Link
      to="/sources/forums"
      style={{ textDecoration: 'none', color: 'inherit' }}
    >
      <ListItemButton sx={{ pl: 4 }}>
        <ListItemIcon style={{ color: '#fff', fontSize: '1.5rem' }}>
          <ForumOutlinedIcon />
        </ListItemIcon>
        <ListItemText
          primary="Forums"
          primaryTypographyProps={{ marginLeft: '-15px' }}
        />
      </ListItemButton>
    </Link>
  </div>
);

export default SourcesForumsItem;
