import React from 'react';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import RadioButtonCheckedOutlinedIcon from '@mui/icons-material/RadioButtonCheckedOutlined';
import { Link } from 'react-router-dom';
import { Divider } from '@mui/material';

const ActivitiesRansomLiveItem = () => (
  <div>
    <Link
      to="/ransomware-monitoring/ransom-live"
      style={{ textDecoration: 'none', color: 'inherit' }}
    >
      <ListItemButton sx={{ pl: 4 }}>
        <ListItemIcon style={{ color: '#fff' }}>
          <RadioButtonCheckedOutlinedIcon />
        </ListItemIcon>
        <ListItemText
          primary="Ransom Live"
          primaryTypographyProps={{ marginLeft: '-15px' }}
        />
      </ListItemButton>
      <Divider />
    </Link>
  </div>
);

export default ActivitiesRansomLiveItem;
