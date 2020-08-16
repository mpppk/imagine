import { Container } from '@material-ui/core';
import CssBaseline from '@material-ui/core/CssBaseline/CssBaseline';
import ThemeProvider from '@material-ui/styles/ThemeProvider';
import { AppProps } from 'next/app';
import React, {FC, useEffect} from 'react';
import { useDispatch, useSelector } from 'react-redux';
import {workspaceActionCreators} from '../actions/workspace';
import {MyAppBar} from '../components/AppBar';
import { State } from '../reducers/reducer';
import { wrapper } from '../store';
import theme from '../theme';

const useWorkSpaceInitializer = () => {
  const dispatch = useDispatch();
  const currentWorkSpace = useSelector((s:State) => s.global.currentWorkSpace)
  useEffect(() => {
    if (currentWorkSpace === null) {
      dispatch(workspaceActionCreators.requestWorkSpaces())
    }
  }, [])
}

// tslint:disable-next-line variable-name
const WrappedApp: FC<AppProps> = ({Component, pageProps}) => {
  useWorkSpaceInitializer()

  return (
    <ThemeProvider theme={theme}>
      {/* CssBaseline kickstart an elegant, consistent, and simple baseline to build upon. */}
      <CssBaseline />
      <MyAppBar/>
      <Container>
        <Component {...pageProps} />
      </Container>
    </ThemeProvider>
  );
}

export default wrapper.withRedux(WrappedApp)
