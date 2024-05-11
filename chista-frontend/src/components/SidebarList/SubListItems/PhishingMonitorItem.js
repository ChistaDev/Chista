import React from 'react';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import ScreenshotMonitorOutlinedIcon from '@mui/icons-material/ScreenshotMonitorOutlined';
import { Link } from 'react-router-dom';
import { Divider } from '@mui/material';

const PhishingMonitorItem = () => (
  <div>
    <Link
      to="/phishing/monitor"
      style={{ textDecoration: 'none', color: 'inherit' }}
    >
      <ListItemButton sx={{ pl: 4 }}>
        <ListItemIcon style={{ color: '#fff' }}>
          <ScreenshotMonitorOutlinedIcon />
        </ListItemIcon>
        <ListItemText
          primary="Monitor"
          primaryTypographyProps={{ marginLeft: '-14px' }}
        />
      </ListItemButton>
      <Divider />
    </Link>
  </div>
);

export default PhishingMonitorItem;
