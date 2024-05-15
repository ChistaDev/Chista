import * as React from 'react';
import Snackbar from '@mui/material/Snackbar';
import Alert from '@mui/material/Alert';
import { useToastMessage } from '../../contexts/ToastMessageContext';

export default function ToastMessage() {
  const { openToast, setOpenToast, severity, toastContent } = useToastMessage();
  const handleClose = (event, reason) => {
    if (reason === 'clickaway') {
      return;
    }

    setOpenToast(false);
  };

  return (
    <div>
      <Snackbar
        open={openToast}
        autoHideDuration={severity === 'warning' ? null : 12000}
        onClose={handleClose}
        anchorOrigin={{
          vertical: 'bottom',
          horizontal: 'right',
        }}
      >
        <Alert
          onClose={handleClose}
          severity={severity}
          variant="filled"
          sx={{ width: '100%' }}
        >
          {toastContent}
        </Alert>
      </Snackbar>
    </div>
  );
}
