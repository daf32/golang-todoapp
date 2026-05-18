import React, { createContext, useContext, useState, useCallback, useEffect } from 'react';
import { ThemeProvider } from './context/ThemeContext';
import LoadingPage    from './pages/LoadingPage';
import OnboardingPage from './pages/OnboardingPage';
import SignUpPage     from './pages/SignUpPage';
import LoginPage      from './pages/LoginPage';
import TasksPage      from './pages/TasksPage';
import CalendarPage   from './pages/CalendarPage';
import SettingsPage   from './pages/SettingsPage';
import ProfilePage    from './pages/ProfilePage';
import './App.css';

export const AuthContext = createContext(null);
export const NavContext  = createContext(null);

export function useAuth() { return useContext(AuthContext); }
export function useNav()  { return useContext(NavContext);  }

const LOADING_MS = 1600;

function AppInner() {
  const [booting, setBooting] = useState(true);
  const [screen,  setScreen]  = useState(() =>
    localStorage.getItem('access_token') ? 'tasks' : 'onboarding'
  );

  /* Splash screen timer */
  useEffect(() => {
    const t = setTimeout(() => setBooting(false), LOADING_MS);
    return () => clearTimeout(t);
  }, []);

  const navigate = useCallback((page) => setScreen(page), []);

  const login = useCallback((accessToken, refreshToken) => {
    localStorage.setItem('access_token', accessToken);
    if (refreshToken) localStorage.setItem('refresh_token', refreshToken);
    setScreen('tasks');
  }, []);

  const logout = useCallback(() => {
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
    setScreen('login');
  }, []);

  const token = localStorage.getItem('access_token');

  const SCREENS = {
    onboarding: <OnboardingPage />,
    signup:     <SignUpPage />,
    login:      <LoginPage />,
    tasks:      <TasksPage />,
    calendar:   <CalendarPage />,
    settings:   <SettingsPage />,
    profile:    <ProfilePage />,
  };

  return (
    <AuthContext.Provider value={{ token, login, logout }}>
      <NavContext.Provider value={{ screen, navigate }}>
        <div className="app-shell">
          <div className="app-phone">
            {booting
              ? <LoadingPage />
              : (
                /* Key forces a re-mount (and therefore re-animation) on page change */
                <div key={screen} className="screen-enter" style={{ flex: 1, display: 'flex', flexDirection: 'column' }}>
                  {SCREENS[screen] ?? <OnboardingPage />}
                </div>
              )
            }
          </div>
        </div>
      </NavContext.Provider>
    </AuthContext.Provider>
  );
}

export default function App() {
  return (
    <ThemeProvider>
      <AppInner />
    </ThemeProvider>
  );
}
