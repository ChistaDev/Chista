import React from 'react';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import GroupsOutlinedIcon from '@mui/icons-material/GroupsOutlined';
import { Link } from 'react-router-dom';

const ThreatProfileAptGroupsItem = () => (
  <div>
    <Link
      to="/threat-profile/apt-groups"
      style={{ textDecoration: 'none', color: 'inherit' }}
    >
      <ListItemButton sx={{ pl: 4 }}>
        <ListItemIcon style={{ color: '#fff' }}>
          <GroupsOutlinedIcon />
        </ListItemIcon>
        <ListItemText
          primary="APT Groups"
          primaryTypographyProps={{
            marginLeft: '-15px',
          }}
        />
      </ListItemButton>
    </Link>
  </div>
);

export default ThreatProfileAptGroupsItem;
