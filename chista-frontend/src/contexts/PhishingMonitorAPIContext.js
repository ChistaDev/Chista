import React, { createContext, useState, useEffect, useContext } from 'react';
import axios from 'axios';
import { useBackendStatus } from './BackendStatusContext';

export const PhishingMonitorAPIContext = createContext();

export const PhishingMonitorAPIProvider = ({ children }) => {
  const [phishingMonitorDetailsData, setPhishingMonitorDetailsData] = useState(
    []
  );
  const { setBackendStatus } = useBackendStatus();

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await axios.get(
          'http://localhost:7777/api/v1/phishing/monitor'
        );
        setBackendStatus(true);
        setPhishingMonitorDetailsData(response.data.results);
        console.log(response.data.results);
      } catch (error) {
        console.error('Error fetching data:', error);
        setBackendStatus(false);
      }
    };

    fetchData();

    const intervalId = setInterval(fetchData, 86400000); // fetch data every 24 hours

    return () => clearInterval(intervalId);
  }, []);

  return (
    <PhishingMonitorAPIContext.Provider
      value={{ phishingMonitorDetailsData, setPhishingMonitorDetailsData }}
    >
      {children}
    </PhishingMonitorAPIContext.Provider>
  );
};

export const usePhishingMonitorAPI = () =>
  useContext(PhishingMonitorAPIContext);
