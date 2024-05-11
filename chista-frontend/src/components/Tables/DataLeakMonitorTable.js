import * as React from 'react';
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
import PhishingScanTable from './PhishingScanTable';

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

const DataLeakMonitorTable = ({
  dataLeakMonitorTableData,
  setDataLeakMonitorTableData,
}) => {
  const [openRowIndex, setOpenRowIndex] = React.useState(null);
  const [openDeleteDialog, setOpenDeleteDialog] = React.useState(false);
  const [deleteIndex, setDeleteIndex] = React.useState(null);
  const [deleteButtonClicked, setDeleteButtonClicked] = React.useState(false);

  const handleModalOpen = (rowIndex) => setOpenRowIndex(rowIndex);

  const handleModalClose = () => setOpenRowIndex(null);

  const handleDelete = (index) => {
    setDeleteIndex(index);
    setOpenDeleteDialog(true);
    setDeleteButtonClicked(false);
  };

  const handleDeleteConfirm = () => {
    const updatedData = [...dataLeakMonitorTableData];
    updatedData.splice(deleteIndex, 1);

    setDataLeakMonitorTableData(updatedData);
    setOpenDeleteDialog(false);
    setDeleteButtonClicked(true);
  };

  const handleDeleteCancel = () => {
    setOpenDeleteDialog(false);
  };

  return (
    <>
      <TableContainer component={Paper}>
        <Table sx={{ minWidth: 650 }} aria-label="caption table">
          <caption>Monitor Data Leak</caption>
          <TableHead>
            <TableRow
              sx={{
                backgroundColor: '#f4f4f4',
              }}
            >
              <TableCell>SN</TableCell>
              <TableCell align="right">Domains</TableCell>
              <TableCell align="right">Created At</TableCell>
              <TableCell align="right">Details</TableCell>
              <TableCell align="right">Delete</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {dataLeakMonitorTableData.map((row, index) => (
              <TableRow key={row.id}>
                <TableCell component="th" scope="row">
                  {row.sn}
                </TableCell>
                <TableCell align="right">
                  {row.dataLeakMonitorDomainInput}
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
                      >
                        {row.dataLeakMonitorDomainInput}
                      </Typography>
                      <Typography id="modal-modal-description" sx={{ mt: 2 }}>
                        <span style={{ fontWeight: 'bold' }}>Created At: </span>
                        {row.createdAt}
                      </Typography>
                      <PhishingScanTable />
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

export default DataLeakMonitorTable;
