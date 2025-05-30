import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { ThemeProvider, createTheme } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import { Provider } from 'react-redux';
import store from './store';

// Components
import Navbar from './components/Navbar';
import PrivateRoute from './components/PrivateRoute';

// Pages
import Home from './pages/Home';
import Login from './pages/Login';
import Register from './pages/Register';
import Cottages from './pages/Cottages';
import Booking from './pages/Booking';
import Profile from './pages/Profile';

const theme = createTheme({
  palette: {
    primary: {
      main: '#4CAF50',
      light: '#81C784',
      dark: '#388E3C',
    },
    secondary: {
      main: '#FFC107',
      light: '#FFE082',
      dark: '#FFA000',
    },
    background: {
      default: '#F5F5F5',
    },
  },
});

function App() {
  return (
    <Provider store={store}>
      <ThemeProvider theme={theme}>
        <CssBaseline />
        <Router>
          <Navbar />
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/login" element={<Login />} />
            <Route path="/register" element={<Register />} />
            <Route path="/cottages" element={<PrivateRoute><Cottages /></PrivateRoute>} />
            <Route path="/booking/:id" element={<PrivateRoute><Booking /></PrivateRoute>} />
            <Route path="/profile" element={<PrivateRoute><Profile /></PrivateRoute>} />
          </Routes>
        </Router>
      </ThemeProvider>
    </Provider>
  );
}

export default App;
