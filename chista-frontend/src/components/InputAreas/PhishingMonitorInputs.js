import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { TextField, Button, Box } from '@mui/material';
import SaveAltOutlinedIcon from '@mui/icons-material/SaveAltOutlined';
import { usePhishingMonitorDomainInput } from '../../contexts/PhishingMonitorDomainInputContext';
import { usePhishingMonitorExcludeInput } from '../../contexts/PhishingMonitorExcludeInputContext';
import PhishingMonitorTable from '../Tables/PhishingMonitorTable';
import ToastMessage from '../ToastMessage/ToastMessage';

const PhishingMonitorInputs = () => {
  const [openToast, setOpenToast] = useState(false);
  const [severity, setSeverity] = useState('');
  const [toastContent, setToastContent] = useState('');

  const [isButtonClicked, setIsButtonClicked] = useState(false);
  const [showDomainError, setShowDomainError] = useState(false);
  const [showExcludeError, setShowExcludeError] = useState(false);
  const { phishingMonitorDomainInput, setPhishingMonitorDomainInput } =
    usePhishingMonitorDomainInput();
  const { phishingMonitorExcludeInput, setPhishingMonitorExcludeInput } =
    usePhishingMonitorExcludeInput();
  const [phishingMonitorTableData, setPhishingMonitorTableData] = useState([]);

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

  const handleSaveButtonClick = () => {
    // Check if any domain is already present in the table
    if (
      phishingMonitorTableData.some(
        (row) =>
          row.phishingMonitorDomainInput === phishingMonitorDomainInput &&
          row.phishingMonitorExcludeInput === phishingMonitorExcludeInput
      )
    ) {
      setSeverity('error');
      setToastContent('This domain is already added to the table.');
      setOpenToast(true);
      return;
    }

    // Validate phishingMonitorDomainInput format
    const domainRegex = /^(?!www\.)[a-zA-Z0-9-]+(\.[a-zA-Z]{2,}){1}$/;
    if (!domainRegex.test(phishingMonitorDomainInput)) {
      setShowDomainError(true);
      return;
    }

    // Validate phishingScanExcludedInput format
    if (!validateExcludeInput(phishingMonitorExcludeInput)) {
      setShowExcludeError(true);
      return;
    }

    if (phishingMonitorDomainInput.trim() !== '') {
      const newRow = {
        sn:
          phishingMonitorTableData.length === 0
            ? 1
            : phishingMonitorTableData[phishingMonitorTableData.length - 1].sn +
              1,
        phishingMonitorDomainInput,
        phishingMonitorExcludeInput,
        createdAt: new Date().toLocaleString(),
      };

      setPhishingMonitorTableData([...phishingMonitorTableData, newRow]);
      setIsButtonClicked(true);
      setShowDomainError(false);
    }
  };

  useEffect(() => {
    if (isButtonClicked) {
      axios
        .post(`http://localhost:7777/api/v1/phishing/monitor`, {
          domain: phishingMonitorDomainInput,
          excludedInput: phishingMonitorExcludeInput,
        })
        .then((response) => {
          // console.log(response);
          // console.log(response.data.msg);
          setIsButtonClicked(false);
          setOpenToast(true);
          if (response.data.msg) {
            setSeverity('success');
            setToastContent(response.data.msg);
          }
          if (response.data.error) {
            setSeverity('error');
            setToastContent(response.data.error);
          }
        })
        .catch((error) => {
          console.error('Error fetching data: ', error);
        });
    }
  }, [isButtonClicked]);

  useEffect(() => {
    // Bileşen ilk kez render edildiğinde localStorage'dan kayıtlı bir tablo var mı kontrol ediyoruz
    const savedTableData = localStorage.getItem('phishingMonitorTableData');
    if (savedTableData) {
      setPhishingMonitorTableData(JSON.parse(savedTableData));
    }
  }, []);

  useEffect(() => {
    // Tablo güncellendiğinde, güncel tabloyu localStorage'a kaydediyoruz
    localStorage.setItem(
      'phishingMonitorTableData',
      JSON.stringify(phishingMonitorTableData)
    );
  }, [phishingMonitorTableData]);

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
          id="domain"
          label="Domain"
          variant="outlined"
          size="small"
          sx={{ width: '650px' }}
          value={phishingMonitorDomainInput}
          onChange={(e) => setPhishingMonitorDomainInput(e.target.value)}
          error={showDomainError}
          helperText={showDomainError ? 'Please enter a valid domain.' : ''}
        />

        <TextField
          id="excluded-input"
          disabled
          label="Exclude Domain (COMING SOON)"
          variant="outlined"
          size="small"
          sx={{ width: '650px' }}
          value={phishingMonitorExcludeInput}
          onChange={(e) =>
            setPhishingMonitorExcludeInput(e.target.value.trim())
          }
          error={showExcludeError}
          helperText={
            showExcludeError
              ? 'Separate the domains using comma (Example: google.com,github.com)'
              : ''
          }
        />

        <Button
          variant="contained"
          onClick={handleSaveButtonClick}
          sx={{
            minWidth: '90px',
            height: '36px',
            display: 'flex',
          }}
          endIcon={<SaveAltOutlinedIcon />}
          disabled={isButtonClicked}
        >
          {isButtonClicked ? <span>SAVING...</span> : <span>SAVE</span>}
        </Button>
      </Box>
      <PhishingMonitorTable
        phishingMonitorTableData={phishingMonitorTableData}
        setPhishingMonitorTableData={setPhishingMonitorTableData}
        setToastContent={setToastContent}
        setSeverity={setSeverity}
        setOpenToast={setOpenToast}
      />
      <ToastMessage
        openToast={openToast}
        setOpenToast={setOpenToast}
        severity={severity}
        setSeverity={setSeverity}
        toastContent={toastContent}
      />
    </>
  );
};

export default PhishingMonitorInputs;
