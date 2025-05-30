import React from 'react';
import { useNavigate } from 'react-router-dom';
import { useSelector } from 'react-redux';
import {
  AppBar,
  Toolbar,
  Typography,
  Button,
  IconButton,
  Box,
  Container,
} from '@mui/material';
import { Person, ExitToApp } from '@mui/icons-material';

const Navbar = () => {
  const navigate = useNavigate();
  const isAuthenticated = useSelector((state) => state.auth.isAuthenticated);

  const handleLogout = () => {
    localStorage.removeItem('token');
    navigate('/login');
  };

  return (
    <AppBar position="static">
      <Container maxWidth="xl">
        <Toolbar>
          <Typography
            variant="h6"
            component="div"
            sx={{ flexGrow: 1, cursor: 'pointer' }}
            onClick={() => navigate('/')}
          >
            Лесная База Отдыха
          </Typography>
          <Box>
            {!isAuthenticated ? (
              <>
                <Button color="inherit" onClick={() => navigate('/login')}>
                  Войти
                </Button>
                <Button color="inherit" onClick={() => navigate('/register')}>
                  Регистрация
                </Button>
              </>
            ) : (
              <>
                <IconButton color="inherit" onClick={() => navigate('/profile')}>
                  <Person />
                </IconButton>
                <IconButton color="inherit" onClick={handleLogout}>
                  <ExitToApp />
                </IconButton>
              </>
            )}
          </Box>
        </Toolbar>
      </Container>
    </AppBar>
  );
};

export default Navbar;
