import React from 'react';
import List from '@mui/material/List';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import ExpandLess from '@mui/icons-material/ExpandLess';
import ExpandMore from '@mui/icons-material/ExpandMore';
import Collapse from '@mui/material/Collapse';
import SourceOutlinedIcon from '@mui/icons-material/SourceOutlined';
import SourcesAptGroupsItem from '../SubListItems/SourcesAptGroupsItem';
import SourcesTelegramItem from '../SubListItems/SourcesTelegramItem';
import SourcesDiscordItem from '../SubListItems/SourcesDiscordItem';
import SourcesBlackMarketsItem from '../SubListItems/SourcesBlackMarketsItem';
import SourcesForumsItem from '../SubListItems/SourcesForumsItem';
import SourcesExploitsItem from '../SubListItems/SourcesExploitsItem';

const SourcesItem = () => {
  const [open, setOpen] = React.useState(false);

  const handleClick = () => {
    setOpen(!open);
  };
  return (
    <div>
      <ListItemButton onClick={handleClick}>
        <ListItemIcon style={{ color: '#fff' }}>
          <SourceOutlinedIcon />
        </ListItemIcon>
        <ListItemText
          primary="Sources"
          primaryTypographyProps={{ marginLeft: '-7.8px' }}
        />
        {open ? <ExpandLess /> : <ExpandMore />}
      </ListItemButton>
      <Collapse in={open} timeout="auto" unmountOnExit>
        <List component="div" disablePadding>
          <SourcesAptGroupsItem />
          <SourcesTelegramItem />
          <SourcesDiscordItem />
          <SourcesBlackMarketsItem />
          <SourcesForumsItem />
          <SourcesExploitsItem />
        </List>
      </Collapse>
    </div>
  );
};

export default SourcesItem;
