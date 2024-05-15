import React from 'react';
import { Typography, Modal, Paper, Button } from '@mui/material';
import CloseOutlinedIcon from '@mui/icons-material/CloseOutlined';

const ModulesModal = ({ isModalOpen, closeModal }) => (
  <Modal
    open={isModalOpen}
    onClose={closeModal}
    aria-labelledby="modal-modal-title"
    aria-describedby="modal-modal-description"
  >
    <Paper
      sx={{
        position: 'absolute',
        top: '50%',
        left: '50%',
        transform: 'translate(-50%, -50%)',
        bgcolor: 'background.paper',
        boxShadow: 24,
        p: 4,
        maxWidth: 800,
      }}
    >
      <Button
        onClick={closeModal}
        style={{ position: 'absolute', top: 0, right: 0 }}
      >
        <CloseOutlinedIcon />
      </Button>
      <Typography
        id="modal-modal-title"
        variant="h4"
        component="h2"
        textAlign={'center'}
      >
        DISCOVER THE MODULES
      </Typography>
      <Typography id="modal-modal-description" sx={{ mt: 2 }}>
        <Typography variant="body1" component="p" gutterBottom mb={3}>
          <b>Collecting IOCs:</b> IOCs are signatures used to identify and track
          cyber threats. Chista can collect IOCs from various sources and make
          them available to users.
        </Typography>
        <Typography variant="body1" component="p" gutterBottom mb={3}>
          <b>Monitoring Data Leaks:</b> Chista can identify accounts that have
          suffered a data breach by monitoring data leaks from various sources.
        </Typography>
        <Typography variant="body1" component="p" gutterBottom mb={3}>
          <b>Monitoring Phishing Campaigns:</b> Chista detects websites created
          for phishing purposes and provides users with a feed in this
          direction.
        </Typography>
        <Typography variant="body1" component="p" gutterBottom mb={3}>
          <b>Monitoring Threat Groups Activities:</b> Threat groups are
          organized groups that carry out cyber attacks. By monitoring threat
          group activity from various sources, Chista helps organizations
          understand and prepare for the activities of threat groups.
        </Typography>
        <Typography variant="body1" component="p" gutterBottom mb={3}>
          <b>Detection of Blacklisted IPs:</b> Chista provides users with a feed
          for IPs blacklisted by various lists.
        </Typography>
        <Typography variant="body1" component="p" gutterBottom mb={3}>
          <b>Detailing Threat Groups:</b> Chista brings together details of
          cybercrime groups, allowing users to get detailed information about
          threat groups.
        </Typography>
        <Typography variant="body1" component="p" gutterBottom mb={3}>
          <b>Detailing Threat Groups:</b> Chista provides resources that can be
          used for threat intelligence for the benefit of users interested in
          Cyber Threat Intelligence.
        </Typography>
      </Typography>
    </Paper>
  </Modal>
);

export default ModulesModal;
