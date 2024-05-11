import React, { useState, useEffect } from 'react';
import TextField from '@mui/material/TextField';
import Button from '@mui/material/Button';
import Box from '@mui/material/Box';
import TravelExploreOutlinedIcon from '@mui/icons-material/TravelExploreOutlined';
import { Fade } from '@mui/material';
import { useDataLeakScanDomainInput } from '../../contexts/DataLeakScanDomainInputContext';
import PhishingScanTable from '../Tables/PhishingScanTable';

const DataLeakScanInnputs = () => {
  const [isButtonClicked, setIsButtonClicked] = useState(false);
  const [displayTable, setDisplayTable] = useState(false);
  const { dataLeakScanDomainInput, setDataLeakScanDomainInput } =
    useDataLeakScanDomainInput();
  const [showDomainError, setShowDomainError] = useState(false);

  const handleButtonClick = () => {
    if (dataLeakScanDomainInput.trim() === '') {
      setShowDomainError(true);
    }

    if (dataLeakScanDomainInput.trim() !== '') {
      setDisplayTable(false);
      setIsButtonClicked(true);
      setShowDomainError(false);
    }
  };

  useEffect(() => {
    if (isButtonClicked) {
      const timeoutId = setTimeout(() => {
        setDisplayTable(true);
        setIsButtonClicked(false);
      }, 1000);

      return () => clearTimeout(timeoutId);
    }
  }, [isButtonClicked]);

  return (
    <>
      <Box
        sx={{
          display: 'flex',
          justifyContent: 'center',
          alignItems: 'top',
          gap: '16px',
          paddingBottom: '20px',
        }}
      >
        <TextField
          id="domain-input"
          label="Domain"
          variant="outlined"
          size="small"
          sx={{ width: '650px' }}
          value={dataLeakScanDomainInput}
          onChange={(e) => setDataLeakScanDomainInput(e.target.value)}
          error={showDomainError}
          helperText={showDomainError ? 'Please enter a valid domain.' : ''}
        />

        <Button
          variant="contained"
          onClick={handleButtonClick}
          sx={{
            minWidth: '90px',
            height: '36px',
          }}
          endIcon={<TravelExploreOutlinedIcon />}
          disabled={isButtonClicked}
        >
          {isButtonClicked ? <span>SCANNING...</span> : <span>SCAN</span>}
        </Button>
      </Box>
      <Fade in={displayTable} timeout={500}>
        <Box>{displayTable ? <PhishingScanTable /> : null}</Box>
      </Fade>
    </>
  );
};

export default DataLeakScanInnputs;
