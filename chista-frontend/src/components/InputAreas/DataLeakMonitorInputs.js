import React, { useState, useEffect } from 'react';
import TextField from '@mui/material/TextField';
import Button from '@mui/material/Button';
import Box from '@mui/material/Box';
import SaveAltOutlinedIcon from '@mui/icons-material/SaveAltOutlined';
import { useDataLeakMonitorDomainInput } from '../../contexts/DataLeakMonitorDomainInputContext';
import DataLeakMonitorTable from '../Tables/DataLeakMonitorTable';

const DataLeakMonitorInputs = () => {
  const [isButtonClicked, setIsButtonClicked] = useState(false);
  const [showDomainError, setShowDomainError] = useState(false);
  const { dataLeakMonitorDomainInput, setDataLeakMonitorDomainInput } =
    useDataLeakMonitorDomainInput();
  const [dataLeakMonitorTableData, setDataLeakMonitorTableData] = useState([
    {
      sn: 1,
      dataLeakMonitorDomainInput: 'dataleak.com',
      createdAt: new Date().toLocaleString(),
    },
  ]);

  const handleButtonClick = () => {
    if (dataLeakMonitorDomainInput.trim() === '') {
      setShowDomainError(true);
    }

    if (dataLeakMonitorDomainInput.trim() !== '') {
      const newRow = {
        sn:
          dataLeakMonitorTableData.length === 0
            ? 1
            : dataLeakMonitorTableData[dataLeakMonitorTableData.length - 1].sn +
              1,
        dataLeakMonitorDomainInput,
        createdAt: new Date().toLocaleString(),
      };

      setDataLeakMonitorTableData([...dataLeakMonitorTableData, newRow]);
      setIsButtonClicked(true);
      setShowDomainError(false);
    }
  };

  useEffect(() => {
    if (isButtonClicked) {
      const timeoutId = setTimeout(() => {
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
          value={dataLeakMonitorDomainInput}
          onChange={(e) => setDataLeakMonitorDomainInput(e.target.value)}
          error={showDomainError}
          helperText={showDomainError ? 'Please enter a valid domain.' : ''}
        />

        <Button
          variant="contained"
          onClick={handleButtonClick}
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
      <DataLeakMonitorTable
        dataLeakMonitorTableData={dataLeakMonitorTableData}
        setDataLeakMonitorTableData={setDataLeakMonitorTableData}
      />
    </>
  );
};

export default DataLeakMonitorInputs;
