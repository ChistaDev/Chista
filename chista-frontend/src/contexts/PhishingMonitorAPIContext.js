import React, { createContext, useState, useEffect, useContext } from 'react';
import axios from 'axios';

export const PhishingMonitorAPIContext = createContext();

export const PhishingMonitorAPIProvider = ({ children }) => {
  const [phishingMonitorData, setPhishingMonitorData] = useState(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await axios.get(
          'http://localhost:7777/api/v1/phishing/monitor'
        );
        setPhishingMonitorData(response.data.results);
        console.log(response.data.results);
      } catch (error) {
        console.error('Error fetching data:', error);
      }
    };

    fetchData();

    const intervalId = setInterval(fetchData, 86400000); // fetch data every 24 hours

    return () => clearInterval(intervalId);
  }, []);

  return (
    <PhishingMonitorAPIContext.Provider value={{ phishingMonitorData }}>
      {children}
    </PhishingMonitorAPIContext.Provider>
  );
};

export const usePhishingMonitorAPI = () =>
  useContext(PhishingMonitorAPIContext);
