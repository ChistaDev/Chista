import React from 'react';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import { TbShoppingCartExclamation } from 'react-icons/tb';
import { Link } from 'react-router-dom';

const SourcesBlackMarketsItem = () => (
  <div>
    <Link
      to="/sources/black-markets"
      style={{ textDecoration: 'none', color: 'inherit' }}
    >
      <ListItemButton sx={{ pl: 4 }}>
        <ListItemIcon style={{ color: '#fff', fontSize: '1.5rem' }}>
          <TbShoppingCartExclamation />
        </ListItemIcon>
        <ListItemText
          primary="Black Markets"
          primaryTypographyProps={{ marginLeft: '-15px' }}
        />
      </ListItemButton>
    </Link>
  </div>
);

export default SourcesBlackMarketsItem;
