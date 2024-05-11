import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { Button, Box, TextField, Fade } from '@mui/material';
import TravelExploreOutlinedIcon from '@mui/icons-material/TravelExploreOutlined';
import { usePhishingScanExcludeInput } from '../../contexts/PhishingScanExcludeInputContext';
import { usePhishingScanDomainInput } from '../../contexts/PhishingScanDomainInputContext';
import PhishingScanTable from '../Tables/PhishingScanTable';
import Loading from '../Loading/Loading';

const PhishingScanInputs = () => {
  const [isButtonClicked, setIsButtonClicked] = useState(false);
  const [displayTable, setDisplayTable] = useState(false);
  const { phishingScanDomainInput, setPhishingScanDomainInput } =
    usePhishingScanDomainInput();
  const { phishingScanExcludeInput, setPhishingScanExcludeInput } =
    usePhishingScanExcludeInput();
  const [showDomainError, setShowDomainError] = useState(false);
  const [showExcludeError, setShowExcludeError] = useState(false);
  const [scanData, setScanData] = useState([]);

  function validateExcludeInput(input) {
    if (input === '') {
      return true;
    }

    const domainRegex = /^(?!www\.)[a-zA-Z0-9-]+(\.[a-zA-Z]{2,}){1}$/;
    const values = input.split(',');

    // eslint-disable-next-line no-plusplus
    for (let i = 0; i < values.length; i++) {
      if (!domainRegex.test(values[i])) {
        return false; // if exclude input not in the right format
      }
    }

    return true; // if exclude input in the right format
  }

  const handleButtonClick = () => {
    // Validate phishingScanDomainInput format
    const domainRegex = /^(?!www\.)[a-zA-Z0-9-]+(\.[a-zA-Z]{2,}){1}$/;
    if (!domainRegex.test(phishingScanDomainInput)) {
      setShowDomainError(true);
      return;
    }

    // Validate phishingScanExcludedInput format
    if (!validateExcludeInput(phishingScanExcludeInput)) {
      setShowExcludeError(true);
      return;
    }

    setDisplayTable(false);
    setIsButtonClicked(true);
    setShowDomainError(false);
    setShowExcludeError(false);
  };

  useEffect(() => {
    if (isButtonClicked) {
      axios
        .get(
          `http://localhost:7777/api/v1/phishing?domain=${phishingScanDomainInput}&exclude=${phishingScanExcludeInput}&verbosity=1`
        )
        .then((response) => {
          setScanData(response.data.possible_phishing_urls);
          setDisplayTable(true);
          setIsButtonClicked(false);
        })
        .catch((error) => {
          console.error('Error fetching data: ', error);
        });
    }
  }, [isButtonClicked]);

  useEffect(() => {
    const storedData = JSON.parse(
      localStorage.getItem('phishingScanTableData')
    );
    if (storedData) {
      setDisplayTable(true);

      setScanData(storedData);
    }
  }, []);

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
          label="Domain (xyz.com)"
          variant="outlined"
          size="small"
          sx={{ width: '650px' }}
          value={phishingScanDomainInput}
          onChange={(e) => setPhishingScanDomainInput(e.target.value)}
          error={showDomainError}
          helperText={
            showDomainError ? 'Please enter a valid domain. e.g. xyz.com' : ''
          }
        />

        <TextField
          id="exclude-input"
          label="Exclude Domain (Optional)"
          variant="outlined"
          size="small"
          sx={{ width: '650px' }}
          value={phishingScanExcludeInput}
          onChange={(e) => setPhishingScanExcludeInput(e.target.value.trim())}
          error={showExcludeError}
          helperText={
            showExcludeError
              ? 'Please separate the domains using comma (Example: google.com,github.com)'
              : ''
          }
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
      {isButtonClicked ? (
        <Loading isButtonClicked={isButtonClicked} />
      ) : (
        <span></span>
      )}
      <Fade in={displayTable} timeout={500}>
        <Box>
          {displayTable ? (
            <PhishingScanTable scanData={scanData} setScanData={setScanData} />
          ) : null}
        </Box>
      </Fade>
    </>
  );
};

export default PhishingScanInputs;
