import CssBaseline from "@mui/material/CssBaseline";
import { QueryClient, QueryClientProvider } from "react-query";
import { AuthProvider } from "components/AuthProvider/AuthProvider";
import type { FC, PropsWithChildren, ReactNode } from "react";
import { HelmetProvider } from "react-helmet-async";
import { AppRouter } from "./AppRouter";
import { ErrorBoundary } from "./components/ErrorBoundary/ErrorBoundary";
import { GlobalSnackbar } from "./components/GlobalSnackbar/GlobalSnackbar";
import theme from "./theme";
import "./theme/globalFonts";
import {
  StyledEngineProvider,
  ThemeProvider as MuiThemeProvider,
} from "@mui/material/styles";
import { ThemeProvider as EmotionThemeProvider } from "@emotion/react";

const defaultQueryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: false,
      refetchOnWindowFocus: false,
    },
  },
});

export const ThemeProviders: FC<PropsWithChildren> = ({ children }) => {
  return (
    <StyledEngineProvider injectFirst>
      <MuiThemeProvider theme={theme.dark}>
        <EmotionThemeProvider theme={theme.dark}>
          <CssBaseline enableColorScheme />
          {children}
        </EmotionThemeProvider>
      </MuiThemeProvider>
    </StyledEngineProvider>
  );
};

interface AppProvidersProps {
  children: ReactNode;
  queryClient?: QueryClient;
}

export const AppProviders: FC<AppProvidersProps> = ({
  children,
  queryClient = defaultQueryClient,
}) => {
  return (
    <HelmetProvider>
      <ThemeProviders>
        <ErrorBoundary>
          <QueryClientProvider client={queryClient}>
            <AuthProvider>
              {children}
              <GlobalSnackbar />
            </AuthProvider>
          </QueryClientProvider>
        </ErrorBoundary>
      </ThemeProviders>
    </HelmetProvider>
  );
};

export const App: FC = () => {
  return (
    <AppProviders>
      <AppRouter />
    </AppProviders>
  );
};
