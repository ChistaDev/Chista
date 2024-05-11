import React from 'react';
import List from '@mui/material/List';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import ExpandLess from '@mui/icons-material/ExpandLess';
import ExpandMore from '@mui/icons-material/ExpandMore';
import Collapse from '@mui/material/Collapse';
import SystemSecurityUpdateWarningOutlinedIcon from '@mui/icons-material/SystemSecurityUpdateWarningOutlined';
import IocFeedItem from '../SubListItems/IocFeedItem';

const IocItem = () => {
  const [open, setOpen] = React.useState(false);

  const handleClick = () => {
    setOpen(!open);
  };
  return (
    <div>
      <ListItemButton onClick={handleClick}>
        <ListItemIcon style={{ color: '#fff' }}>
          <SystemSecurityUpdateWarningOutlinedIcon />
        </ListItemIcon>
        <ListItemText
          primary="IOC"
          primaryTypographyProps={{ marginLeft: '-8px' }}
        />
        {open ? <ExpandLess /> : <ExpandMore />}
      </ListItemButton>
      <Collapse in={open} timeout="auto" unmountOnExit>
        <List component="div" disablePadding>
          <IocFeedItem />
        </List>
      </Collapse>
    </div>
  );
};

export default IocItem;
