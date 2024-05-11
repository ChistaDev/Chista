import * as React from 'react';
import Snackbar from '@mui/material/Snackbar';
import Alert from '@mui/material/Alert';

export default function ToastMessage({
  openToast,
  setOpenToast,
  severity,
  toastContent,
}) {
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
        autoHideDuration={6000}
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
