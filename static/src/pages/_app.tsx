import { Container } from '@material-ui/core';
import CssBaseline from '@material-ui/core/CssBaseline/CssBaseline';
import ThemeProvider from '@material-ui/styles/ThemeProvider';
import { AppProps } from 'next/app';
import React, { FC } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { globalAsyncActionCreators } from '../actions/global';
import {MyAppBar} from '../components/AppBar';
import { State } from '../reducers/reducer';
import { wrapper } from '../store';
import theme from '../theme';

const useHandlers = () => {
  const dispatch = useDispatch();
  return {
    logout: () => {
      dispatch(globalAsyncActionCreators.signOut.started(undefined));
    },
    empty: () => {} //tslint:disable-line
  };
};

// tslint:disable-next-line variable-name
const WrappedApp: FC<AppProps> = ({Component, pageProps}) => {
  const handlers = useHandlers();
  const user = useSelector((state: State) => {
    return state.global.user
  });
  return (
    <ThemeProvider theme={theme}>
      {/* CssBaseline kickstart an elegant, consistent, and simple baseline to build upon. */}
      <CssBaseline />
      <MyAppBar user={user} onClickLogout={handlers.logout} />
      <Container>
        <Component {...pageProps} />
      </Container>
    </ThemeProvider>
  );
}

export default wrapper.withRedux(WrappedApp)
