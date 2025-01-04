import ReactDOM from 'react-dom/client'
import { RouterProvider } from 'react-router-dom'
import router from './config/Router.tsx'
import './index.scss'
import { ThemeProvider } from '@emotion/react'
import theme from './config/theme.ts'

ReactDOM.createRoot(document.getElementById('root')!).render(
  // <React.StrictMode>
        <ThemeProvider theme={theme}>
          <RouterProvider router={router}/> 
        </ThemeProvider>
  // </React.StrictMode>
)
