import React from 'react';
import List from '@mui/material/List';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import ExpandLess from '@mui/icons-material/ExpandLess';
import ExpandMore from '@mui/icons-material/ExpandMore';
import Collapse from '@mui/material/Collapse';
import GppMaybeOutlinedIcon from '@mui/icons-material/GppMaybeOutlined';
import ThreatProfileAptGroupsItem from '../SubListItems/ThreatProfileAptGroupsItem';
import ThreatProfileRansomwareGroupsItem from '../SubListItems/ThreatProfileRansomwareGroupsItem';

const ThreatProfileItem = () => {
  const [open, setOpen] = React.useState(false);

  const handleClick = () => {
    setOpen(!open);
  };
  return (
    <div>
      <ListItemButton onClick={handleClick}>
        <ListItemIcon style={{ color: '#fff' }}>
          <GppMaybeOutlinedIcon />
        </ListItemIcon>
        <ListItemText
          primary="Threat Profile"
          primaryTypographyProps={{ marginLeft: '-7.3px' }}
        />
        {open ? <ExpandLess /> : <ExpandMore />}
      </ListItemButton>
      <Collapse in={open} timeout="auto" unmountOnExit>
        <List component="div" disablePadding>
          <ThreatProfileAptGroupsItem />
          <ThreatProfileRansomwareGroupsItem />
        </List>
      </Collapse>
    </div>
  );
};

export default ThreatProfileItem;
