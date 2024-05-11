import React from 'react';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import SearchOutlinedIcon from '@mui/icons-material/SearchOutlined';
import { Link } from 'react-router-dom';

const BlackListScanItem = () => (
  <div>
    <Link
      to="/black-list/scan"
      style={{ textDecoration: 'none', color: 'inherit' }}
    >
      <ListItemButton sx={{ pl: 4 }}>
        <ListItemIcon style={{ color: '#fff' }}>
          <SearchOutlinedIcon />
        </ListItemIcon>
        <ListItemText
          primary="Scan"
          primaryTypographyProps={{ marginLeft: '-15px' }}
        />
      </ListItemButton>
    </Link>
  </div>
);

export default BlackListScanItem;
