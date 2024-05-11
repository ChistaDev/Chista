import React from 'react';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import GroupOutlinedIcon from '@mui/icons-material/GroupOutlined';
import { Link } from 'react-router-dom';
import { Divider } from '@mui/material';

const ThreatProfileRansomwareGroupsItem = () => (
  <div>
    <Link
      to="/threat-profile/ransomware-groups"
      style={{ textDecoration: 'none', color: 'inherit' }}
    >
      <ListItemButton sx={{ pl: 4 }}>
        <ListItemIcon style={{ color: '#fff' }}>
          <GroupOutlinedIcon />
        </ListItemIcon>
        <ListItemText
          primary="Ransomware Groups"
          primaryTypographyProps={{
            noWrap: false,
            sx: { fontSize: '0.98rem', marginLeft: '-15px' },
          }}
        />
      </ListItemButton>
      <Divider />
    </Link>
  </div>
);

export default ThreatProfileRansomwareGroupsItem;
