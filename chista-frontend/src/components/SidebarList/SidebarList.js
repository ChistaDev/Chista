import * as React from 'react';
import List from '@mui/material/List';
import { Divider } from '@mui/material';
import HomePageItem from './ListItems/HomePageItem';
// import DataLeakItem from './ListItems/DataLeakItem';
// import ThreatProfileItem from './ListItems/ThreatProfileItem';
// import ActivitiesItem from './ListItems/ActivitiesItem';
// import BlackListItem from './ListItems/BlackListItem';
// import IocItem from './ListItems/IocItem';
// import SourcesItem from './ListItems/SourcesItem';
// import SettingsItem from './ListItems/SettingsItem';
import PhishingItem from './ListItems/PhishingItem';

const SidebarList = () => (
  <List
    sx={{ width: '100%', maxWidth: 360, color: '#fff' }}
    component="nav"
    aria-labelledby="nested-list-subheader"
  >
    <HomePageItem />
    <Divider />
    <PhishingItem />
    {/* <DataLeakItem />
    <ThreatProfileItem />
    <ActivitiesItem />
    <BlackListItem />
    <IocItem />
    <SourcesItem />
    <Divider />
    <SettingsItem /> */}
  </List>
);

export default SidebarList;
