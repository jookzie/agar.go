import { createTheme } from "@mui/material";

const spotifyGreen = '#1DB954';
const spotifyBlack = '#121212';
const spotifyGrey = '#535353';
const spotifyLightGrey = '#b3b3b3';

const theme = createTheme({
  palette: {
    primary: {
      main: spotifyGreen,
      contrastText: '#ffffff',
    },
    secondary: {
      main: spotifyGrey,
    },
    background: {
      default: spotifyBlack,
      paper: spotifyBlack,
    },
    text: {
      primary: '#ffffff',
      secondary: spotifyLightGrey,
    },
  },
  typography: {
    fontFamily: 'Roboto, sans-serif',
  },
  components: {
    MuiSvgIcon: {
      styleOverrides: {
        root: {
          color: '#ffffff',
        },
      },
    },
    MuiOutlinedInput: {
      styleOverrides: {
        root: {
          '& .MuiOutlinedInput-notchedOutline': {
            borderColor: spotifyGrey,
          },
          '&:hover .MuiOutlinedInput-notchedOutline': {
            borderColor: spotifyGreen,
          },
          '&.Mui-focused .MuiOutlinedInput-notchedOutline': {
            borderColor: spotifyGreen,
          },
        },
        input: {
          color: spotifyLightGrey,
        },
      },
    },
    MuiInputLabel: {
      styleOverrides: {
        root: {
          color: spotifyLightGrey,
          '&.Mui-focused': {
            color: spotifyGreen,
          },
        },
      },
    },
    MuiSelect: {
      styleOverrides: {
        icon: {
          color: 'white',
        },
      },
    },
  },
});

export default theme;