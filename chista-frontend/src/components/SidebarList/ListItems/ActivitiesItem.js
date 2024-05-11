import React from 'react';
import List from '@mui/material/List';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import ExpandLess from '@mui/icons-material/ExpandLess';
import ExpandMore from '@mui/icons-material/ExpandMore';
import Collapse from '@mui/material/Collapse';
import TvOutlinedIcon from '@mui/icons-material/TvOutlined';
import ActivitiesRansomLiveItem from '../SubListItems/ActivitiesRansomLiveItem';

const ActivitiesItem = () => {
  const [open, setOpen] = React.useState(false);

  const handleClick = () => {
    setOpen(!open);
  };
  return (
    <div>
      <ListItemButton onClick={handleClick}>
        <ListItemIcon style={{ color: '#fff' }}>
          <TvOutlinedIcon />
        </ListItemIcon>
        <ListItemText
          primary="Activities"
          primaryTypographyProps={{ marginLeft: '-7.3px' }}
        />
        {open ? <ExpandLess /> : <ExpandMore />}
      </ListItemButton>
      <Collapse in={open} timeout="auto" unmountOnExit>
        <List component="div" disablePadding>
          <ActivitiesRansomLiveItem />
        </List>
      </Collapse>
    </div>
  );
};

export default ActivitiesItem;
