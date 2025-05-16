import { ThemeProvider } from "@emotion/react";
import { Container, createTheme, CssBaseline, Typography } from "@mui/material";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { SnackbarProvider } from 'notistack';

const darkTheme = createTheme({
  palette: {
    mode: 'dark',
  },
});

const queryClient = new QueryClient();

export function App() {
  return (
    <ThemeProvider theme={darkTheme}>
      <QueryClientProvider client={queryClient}>
        <CssBaseline />
        <SnackbarProvider />

        <Container style={{ height: '100%', paddingTop: '16px', paddingBottom: '140px', overflow: 'hidden' }}>
          <Typography variant="h4">Manager</Typography>

        </Container>
      </QueryClientProvider>
    </ThemeProvider>
  );
}

