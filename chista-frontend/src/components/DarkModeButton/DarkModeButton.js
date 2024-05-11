import * as React from 'react';
import IconButton from '@mui/material/IconButton';
import SettingsBrightnessOutlinedIcon from '@mui/icons-material/SettingsBrightnessOutlined';
import { useDarkMode } from '../../contexts/DarkModeContext';

const DarkModeButton = () => {
  const { mode, setMode } = useDarkMode();

  React.useEffect(() => {
    const savedMode = localStorage.getItem('darkMode');
    if (savedMode) {
      setMode(savedMode === 'dark');
    }
  }, []);

  const handleToggle = () => {
    setMode((prevMode) => !prevMode);

    localStorage.setItem('darkMode', !mode ? 'dark' : 'light');
  };

  return (
    <IconButton
      onClick={handleToggle}
      sx={{
        backgroundColor: 'transparent',
        color: 'inherit',
        '&:hover': {
          backgroundColor: 'transparent',
        },
      }}
    >
      <SettingsBrightnessOutlinedIcon />
    </IconButton>
  );
};

export default DarkModeButton;
