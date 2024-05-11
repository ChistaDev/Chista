/* eslint-disable no-nested-ternary */
import * as React from 'react';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TablePagination from '@mui/material/TablePagination';
import TableRow from '@mui/material/TableRow';
import { Box } from '@mui/material';
import MOCK_DATA from './MOCK_DATA.json';

const columns = [
  { id: 'domain', label: 'Domain', minWidth: 170 },
  { id: 'ipaddress', label: 'Ip Address', minWidth: 150 },
  {
    id: 'createdate',
    label: 'Create Date',
    minWidth: 170,
    // align: "right", looks better with numbers
  },
  {
    id: 'source',
    label: 'Source',
    minWidth: 170,
  },
  {
    id: 'country',
    label: 'Country',
    minWidth: 120,
  },
  {
    id: 'nameserver',
    label: 'Name Server',
    minWidth: 170,
  },
  {
    id: 'openport',
    label: 'Open Port',
    minWidth: 170,
  },
  {
    id: 'pori',
    label: 'P or I',
    minWidth: 170,
  },
  {
    id: 'url',
    label: 'URL',
    minWidth: 170,
  },
];

const DataLeakScanTable = () => {
  const [page, setPage] = React.useState(0);
  const [rowsPerPage, setRowsPerPage] = React.useState(10);

  const handleChangeRowsPerPage = (event) => {
    setRowsPerPage(+event.target.value);
    setPage(0);
  };

  const handleChangePage = (event, newPage) => {
    setPage(newPage);
  };

  return (
    <Box
      sx={{
        width: '100%',
        overflow: 'hidden',
        display: 'flex',
        flexDirection: 'column',
        justifyContent: 'center',
        alignItems: 'center',
        gap: '16px',
        paddingBottom: '15px',
        marginTop: '12px',
      }}
    >
      <TableContainer sx={{ maxHeight: 440 }}>
        <Table stickyHeader aria-label="sticky table">
          <TableHead>
            <TableRow>
              {columns.map((column) => (
                <TableCell
                  key={column.id}
                  align={column.align}
                  style={{
                    minWidth: column.minWidth,
                    backgroundColor: '#f4f4f4',
                  }}
                >
                  {column.label}
                </TableCell>
              ))}
            </TableRow>
          </TableHead>
          <TableBody>
            {MOCK_DATA.slice(
              page * rowsPerPage,
              page * rowsPerPage + rowsPerPage
            ).map((row) => (
              <TableRow hover role="checkbox" tabIndex={-1} key={row.id}>
                {columns.map((column) => {
                  const value = row[column.id];
                  return (
                    <TableCell key={column.id} align={column.align}>
                      {value === true ? 'P' : value === false ? 'I' : value}
                    </TableCell>
                  );
                })}
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
      <TablePagination
        rowsPerPageOptions={[10, 25, 100]}
        component="div"
        count={MOCK_DATA.length}
        rowsPerPage={rowsPerPage}
        page={page}
        onPageChange={handleChangePage}
        onRowsPerPageChange={handleChangeRowsPerPage}
        style={{
          backgroundColor: '#f4f4f4',
          borderRadius: '8px',
        }}
      />
    </Box>
  );
};

export default DataLeakScanTable;
