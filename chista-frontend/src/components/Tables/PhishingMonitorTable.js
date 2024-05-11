import * as React from 'react';
import axios from 'axios';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Paper from '@mui/material/Paper';
import {
  Button,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
} from '@mui/material';
import MoreHorizOutlinedIcon from '@mui/icons-material/MoreHorizOutlined';
import CloseOutlinedIcon from '@mui/icons-material/CloseOutlined';
import DeleteOutlineOutlinedIcon from '@mui/icons-material/DeleteOutlineOutlined';
import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import Modal from '@mui/material/Modal';
import PhishingMonitorDetailsTable from './PhishingMonitorDetailsTable';
import { useDarkMode } from '../../contexts/DarkModeContext';

const style = {
  position: 'absolute',
  top: '50%',
  left: '50%',
  transform: 'translate(-50%, -50%)',
  width: '90%',
  maxWidth: '1650px',
  maxHeight: '90%',
  overflowY: 'auto',
  bgcolor: 'background.paper',
  border: '2px solid #000',
  boxShadow: 24,
  p: 4,
};

const PhishingMonitorTable = ({
  phishingMonitorTableData,
  setPhishingMonitorTableData,
  setToastContent,
  setSeverity,
  setOpenToast,
}) => {
  const { mode } = useDarkMode();
  const [openRowIndex, setOpenRowIndex] = React.useState(null);
  const [openDeleteDialog, setOpenDeleteDialog] = React.useState(false);
  const [deleteIndex, setDeleteIndex] = React.useState(null);
  const [deleteButtonClicked, setDeleteButtonClicked] = React.useState(false);

  const handleModalOpen = (rowIndex) => {
    setOpenRowIndex(rowIndex);
    console.log('Clicked item index:', rowIndex);
  };

  const handleModalClose = () => setOpenRowIndex(null);

  const handleDelete = (index) => {
    setDeleteIndex(index);
    setOpenDeleteDialog(true);
    setDeleteButtonClicked(false);
  };

  const handleDeleteConfirm = () => {
    const updatedData = [...phishingMonitorTableData];
    updatedData.splice(deleteIndex, 1);

    setPhishingMonitorTableData(updatedData);
    setOpenDeleteDialog(false);
    setDeleteButtonClicked(true);

    axios
      .delete(
        `http://localhost:7777/api/v1/phishing/monitor?domain=${phishingMonitorTableData[deleteIndex].phishingMonitorDomainInput}`
      )
      .then((response) => {
        setToastContent(response.data.msg);
        setSeverity('info');
        setOpenToast(true);
      })
      .catch((error) => {
        console.error('Error fetching data: ', error);
      });
  };

  const handleDeleteCancel = () => {
    setOpenDeleteDialog(false);
  };

  return (
    <>
      <TableContainer component={Paper}>
        <Table sx={{ minWidth: 650 }} aria-label="caption table">
          <caption>Phishing Monitor </caption>
          <TableHead>
            <TableRow
              sx={
                {
                  // backgroundColor: '#f4f4f4',
                }
              }
            >
              <TableCell>SN</TableCell>
              <TableCell align="right">DOMAIN</TableCell>
              <TableCell align="right">CREATED AT</TableCell>
              <TableCell align="right">DETAILS</TableCell>
              <TableCell align="right">DELETE</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {phishingMonitorTableData.map((row, index) => (
              <TableRow key={index}>
                <TableCell component="th" scope="row">
                  {row.sn}
                </TableCell>
                <TableCell align="right">
                  {row.phishingMonitorDomainInput}
                </TableCell>
                <TableCell align="right">{row.createdAt}</TableCell>
                <TableCell align="right">
                  <Button onClick={() => handleModalOpen(index)}>
                    <MoreHorizOutlinedIcon />
                  </Button>
                  <Modal
                    open={openRowIndex === index}
                    onClose={handleModalClose}
                    aria-labelledby="modal-modal-title"
                    aria-describedby="modal-modal-description"
                  >
                    <Box sx={style}>
                      <Button
                        onClick={handleModalClose}
                        style={{ position: 'absolute', top: 0, right: 0 }}
                      >
                        <CloseOutlinedIcon />
                      </Button>
                      <Typography
                        id="modal-modal-title"
                        variant="h3"
                        component="h2"
                        sx={{ color: mode ? '#fff' : 'rgba(0, 0, 0, 0.87)' }}
                      >
                        {row.phishingMonitorDomainInput}
                      </Typography>
                      <Typography
                        sx={{
                          mt: 2,
                          color: mode ? '#fff' : 'rgba(0, 0, 0, 0.87)',
                        }}
                      >
                        <span style={{ fontWeight: 'bold' }}>Created At: </span>
                        {row.createdAt}
                      </Typography>
                      {/* <Typography
                        sx={{
                          mt: 2,
                          color: mode ? '#fff' : 'rgba(0, 0, 0, 0.87)',
                        }}
                      >
                        <span style={{ fontWeight: 'bold' }}>
                          Excluded Domains:
                        </span>{' '}
                        {row.phishingMonitorExcludeInput}
                      </Typography> */}
                      <PhishingMonitorDetailsTable
                        phishingMonitorTableData={phishingMonitorTableData}
                        openRowIndex={openRowIndex}
                      />
                    </Box>
                  </Modal>
                </TableCell>
                <TableCell align="right">
                  <Button
                    sx={{ color: 'red' }}
                    onClick={() => handleDelete(index)}
                  >
                    <DeleteOutlineOutlinedIcon />
                  </Button>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
      <Dialog
        open={openDeleteDialog}
        onClose={handleDeleteCancel}
        aria-labelledby="alert-dialog-title"
        aria-describedby="alert-dialog-description"
      >
        <DialogTitle id="alert-dialog-title">Confirm Delete</DialogTitle>
        <DialogContent>
          <Typography>Are you sure you want to delete this item?</Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleDeleteCancel}>Cancel</Button>
          <Button
            variant="contained"
            color="error"
            onClick={handleDeleteConfirm}
            disabled={deleteButtonClicked}
          >
            Delete
          </Button>
        </DialogActions>
      </Dialog>
    </>
  );
};

export default PhishingMonitorTable;
