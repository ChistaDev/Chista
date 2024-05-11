/* eslint-disable no-nested-ternary */
import * as React from 'react';
import { useEffect } from 'react';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TablePagination from '@mui/material/TablePagination';
import TableRow from '@mui/material/TableRow';
import { Box } from '@mui/material';
import urlParse from 'url-parse';

const columns = [
  {
    id: 'url',
    label: 'URL',
    minWidth: 170,
  },
  { id: 'domain', label: 'DOMAIN', minWidth: 170 },
  {
    id: 'port',
    label: 'PORT',
    minWidth: 170,
  },
];

const PhishingScanTable = ({ scanData }) => {
  const [page, setPage] = React.useState(0);
  const [rowsPerPage, setRowsPerPage] = React.useState(10);

  const handleChangeRowsPerPage = (event) => {
    setRowsPerPage(+event.target.value);
    setPage(0);
  };

  const handleChangePage = (event, newPage) => {
    setPage(newPage);
  };

  const getDomainAndPort = (url) => {
    const parsedUrl = urlParse(url);
    let { port } = parsedUrl;
    if (!port) {
      port = parsedUrl.protocol === 'https:' ? '443' : '80';
    }
    return {
      domain: parsedUrl.hostname,
      port,
    };
  };

  useEffect(() => {
    localStorage.setItem('phishingScanTableData', JSON.stringify(scanData));
  }, [scanData]);

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
          <caption>Phishing Scan </caption>
          <TableHead>
            <TableRow>
              {columns.map((column) => (
                <TableCell
                  key={column.id}
                  align={column.align}
                  style={{
                    minWidth: column.minWidth,
                    // backgroundColor: '#f4f4f4',
                  }}
                >
                  {column.label}
                </TableCell>
              ))}
            </TableRow>
          </TableHead>
          <TableBody>
            {scanData
              .slice(page * rowsPerPage, page * rowsPerPage + rowsPerPage)
              .map((url, index) => {
                const { domain, port } = getDomainAndPort(url);
                return (
                  <TableRow hover role="checkbox" tabIndex={-1} key={index}>
                    <TableCell align="left">{url}</TableCell>
                    <TableCell align="left">{domain}</TableCell>
                    <TableCell align="left">{port}</TableCell>
                  </TableRow>
                );
              })}
          </TableBody>
        </Table>
      </TableContainer>
      <TablePagination
        rowsPerPageOptions={[10, 25, 100]}
        component="div"
        count={scanData.length}
        rowsPerPage={rowsPerPage}
        page={page}
        onPageChange={handleChangePage}
        onRowsPerPageChange={handleChangeRowsPerPage}
        style={{
          // backgroundColor: '#f4f4f4',
          borderRadius: '8px',
        }}
      />
    </Box>
  );
};

export default PhishingScanTable;
