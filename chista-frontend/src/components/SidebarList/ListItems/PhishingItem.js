import React from 'react';
import List from '@mui/material/List';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import ExpandLess from '@mui/icons-material/ExpandLess';
import ExpandMore from '@mui/icons-material/ExpandMore';
import Collapse from '@mui/material/Collapse';
import PhishingOutlinedIcon from '@mui/icons-material/PhishingOutlined';
import PhishingScanItem from '../SubListItems/PhishingScanItem';
import PhishingMonitorItem from '../SubListItems/PhishingMonitorItem';

const PhishingItem = () => {
  const [open, setOpen] = React.useState(false);

  const handleClick = () => {
    setOpen(!open);
  };
  return (
    <div>
      <ListItemButton onClick={handleClick}>
        <ListItemIcon style={{ color: '#fff' }}>
          <PhishingOutlinedIcon />
        </ListItemIcon>
        <ListItemText
          primary="Phishing"
          primaryTypographyProps={{ marginLeft: '-7px' }}
        />
        {open ? <ExpandLess /> : <ExpandMore />}
      </ListItemButton>
      <Collapse in={open} timeout="auto" unmountOnExit>
        <List component="div" disablePadding>
          <PhishingScanItem />
          <PhishingMonitorItem />
        </List>
      </Collapse>
    </div>
  );
};

export default PhishingItem;
