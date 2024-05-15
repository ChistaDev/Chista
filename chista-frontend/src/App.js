import React, { useEffect } from 'react';
import { Route, Routes } from 'react-router-dom';
import { Box, ThemeProvider, createTheme } from '@mui/material';
import HomePage from './pages/HomePage';
import DataLeakScan from './pages/DataLeakScan';
import DataLeakMonitor from './pages/DataLeakMonitor';
import NotFoundPage from './pages/NotFoundPage';
import NavigationLayout from './pages/NavigationLayout';
import { IsDrawerOpenProvider } from './contexts/IsDrawerOpenContext';
import { PhishingScanExcludeInputProvider } from './contexts/PhishingScanExcludeInputContext';
import { PhishingScanDomainInputProvider } from './contexts/PhishingScanDomainInputContext';
import { PhishingMonitorDomainInputProvider } from './contexts/PhishingMonitorDomainInputContext';
import { PhishingMonitorExcludeInputProvider } from './contexts/PhishingMonitorExcludeInputContext';
import { DataLeakScanDomainInputProvider } from './contexts/DataLeakScanDomainInputContext';
import { DataLeakMonitorDomainInputProvider } from './contexts/DataLeakMonitorDomainInputContext';
import { PhishingMonitorAPIProvider } from './contexts/PhishingMonitorAPIContext';
import { useDarkMode } from './contexts/DarkModeContext';
import { useBackendStatus } from './contexts/BackendStatusContext';
import { useToastMessage } from './contexts/ToastMessageContext';
import PhishingScan from './pages/PhishingScan';
import PhishingMonitor from './pages/PhishingMonitor';
import ThreatProfileAptGroups from './pages/ThreatProfileAptGroups';
import ThreatProfileRansomwareGroups from './pages/ThreatProfileRansomwareGroups';
import ActivitiesRansomLive from './pages/ActivitiesRansomLive';
import BlackListScan from './pages/BlackListScan';
import BlackListMonitor from './pages/BlackListMonitor';
import IocFeed from './pages/IocFeed';
import SourcesAptGroups from './pages/SourcesAptGroups';
import SourcesDiscordChannels from './pages/SourcesDiscordChannels';
import SourcesTelegramChannels from './pages/SourcesTelegramChannels';
import SourcesBlackMarkets from './pages/SourcesBlackMarkets';
import SourcesForums from './pages/SourcesForums';
import SourcesExploits from './pages/SourcesExploits';
import SettingsPage from './pages/SettingsPage';

function App() {
  const { mode } = useDarkMode();

  const { backendStatus } = useBackendStatus();

  const { setOpenToast, setSeverity, setToastContent } = useToastMessage();

  const darkTheme = createTheme({
    palette: {
      mode: mode ? 'dark' : 'light',
    },
  });

  useEffect(() => {
    if (!backendStatus) {
      setOpenToast(true);
      setSeverity('warning');
      setToastContent('Please make sure the backend server is running.');
    } else {
      setOpenToast(false);
    }
  }, [backendStatus]);

  return (
    <Box>
      <ThemeProvider theme={darkTheme}>
        <IsDrawerOpenProvider>
          <PhishingScanExcludeInputProvider>
            <PhishingScanDomainInputProvider>
              <PhishingMonitorDomainInputProvider>
                <PhishingMonitorExcludeInputProvider>
                  <DataLeakScanDomainInputProvider>
                    <DataLeakMonitorDomainInputProvider>
                      <PhishingMonitorAPIProvider>
                        <Routes>
                          <Route element={<NavigationLayout />}>
                            <Route path="/" element={<HomePage />} />

                            <Route
                              path="/phishing/scan"
                              element={<PhishingScan />}
                            />

                            <Route
                              path="/phishing/monitor"
                              element={<PhishingMonitor />}
                            />

                            <Route
                              path="/data-leak/scan"
                              element={<DataLeakScan />}
                            />

                            <Route
                              path="/data-leak/monitor"
                              element={<DataLeakMonitor />}
                            />

                            <Route path="/threat-profile">
                              <Route
                                path="apt-groups"
                                element={<ThreatProfileAptGroups />}
                              />
                              <Route
                                path="ransomware-groups"
                                element={<ThreatProfileRansomwareGroups />}
                              />
                            </Route>

                            <Route path="/ransomware-monitoring">
                              <Route
                                path="ransom-live"
                                element={<ActivitiesRansomLive />}
                              />
                            </Route>

                            <Route path="/black-list">
                              <Route path="scan" element={<BlackListScan />} />
                              <Route
                                path="monitor"
                                element={<BlackListMonitor />}
                              />
                            </Route>

                            <Route path="/ioc">
                              <Route path="ioc-feed" element={<IocFeed />} />
                            </Route>

                            <Route path="/sources">
                              <Route
                                path="apt-groups"
                                element={<SourcesAptGroups />}
                              />

                              <Route
                                path="telegram-channels"
                                element={<SourcesTelegramChannels />}
                              />
                              <Route
                                path="discord-channels"
                                element={<SourcesDiscordChannels />}
                              />
                              <Route
                                path="black-markets"
                                element={<SourcesBlackMarkets />}
                              />
                              <Route
                                path="forums"
                                element={<SourcesForums />}
                              />
                              <Route
                                path="exploits"
                                element={<SourcesExploits />}
                              />
                            </Route>

                            <Route
                              path="/settings"
                              element={<SettingsPage />}
                            />

                            <Route path="*" element={<NotFoundPage />} />
                          </Route>
                        </Routes>
                      </PhishingMonitorAPIProvider>
                    </DataLeakMonitorDomainInputProvider>
                  </DataLeakScanDomainInputProvider>
                </PhishingMonitorExcludeInputProvider>
              </PhishingMonitorDomainInputProvider>
            </PhishingScanDomainInputProvider>
          </PhishingScanExcludeInputProvider>
        </IsDrawerOpenProvider>
      </ThemeProvider>
    </Box>
  );
}

export default App;
