import { CircularProgress } from '@material-ui/core';
import Avatar from '@material-ui/core/Avatar';
import Box from '@material-ui/core/Box';
import Button from '@material-ui/core/Button';
import Checkbox from '@material-ui/core/Checkbox';
import { green } from '@material-ui/core/colors';
import Container from '@material-ui/core/Container';
import CssBaseline from '@material-ui/core/CssBaseline';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import Grid from '@material-ui/core/Grid';
import Link from '@material-ui/core/Link';
import { makeStyles } from '@material-ui/core/styles';
import TextField from '@material-ui/core/TextField';
import Typography from '@material-ui/core/Typography';
import LockOutlinedIcon from '@material-ui/icons/LockOutlined';
import { NextPage } from 'next';
import { useRouter } from 'next/router';
import React, { useEffect, useState } from 'react';
import { useSelector } from 'react-redux';
import { globalActionCreators } from '../actions/global';
import { useActions } from '../hooks';
import { State } from '../reducers/reducer';

function Copyright() {
  return (
    <Typography variant="body2" color="textSecondary" align="center">
      {'Copyright © '}
      <Link color="inherit" href="https://github.com/mpppk/next-ts-redux-material/">
        mpppk
      </Link>{' '}
      {new Date().getFullYear()}
      {'.'}
    </Typography>
  );
}

const useStyles = makeStyles((theme) => ({
  avatar: {
    backgroundColor: theme.palette.secondary.main,
    margin: theme.spacing(1),
  },
  buttonProgress: {
    color: green[500],
    left: '50%',
    marginLeft: -12,
    marginTop: -12+4,
    position: 'absolute',
    top: '50%',
  },
  form: {
    marginTop: theme.spacing(1),
    width: '100%', // Fix IE 11 issue.
  },
  paper: {
    alignItems: 'center',
    display: 'flex',
    flexDirection: 'column',
    marginTop: theme.spacing(8),
  },
  submit: {
    margin: theme.spacing(3, 0, 2),
  },
  wrapper: {
    margin: theme.spacing(1),
    position: 'relative',
  },
}));

const useComponentState = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  return {
    email,
    password,
    setEmail,
    setPassword,
  };
};

const useHandlers = () => {
  const componentState = useComponentState();
  const actionCreators = useActions(globalActionCreators);
  return {
    onChangeEmailForm: (e: React.ChangeEvent<HTMLInputElement>) => {
      componentState.setEmail(e.target.value);
    },
    onChangePasswordForm: (e: React.ChangeEvent<HTMLInputElement>) => {
      componentState.setPassword(e.target.value);
    },
    onClickSignInSubmitButton: () => {
      actionCreators.clickSignInSubmitButton({
        email: componentState.email,
        password: componentState.password,
      });
    }
  };
};

const useSignInRouter = () => {
  const isSignedIn = useSelector((s: State) => !!s.global.jwt)
  const router = useRouter();
  useEffect(() => {
    if (isSignedIn) {
      router.push('/');
    }
  }, [isSignedIn])
}


// tslint:disable-next-line variable-name
export const SignIn: NextPage = () => {
  const state = useSelector((s: State) => ({
    signedIn: !!s.global.jwt,
    waitingSignIn: s.global.waitingSignIn,
  }))

  useSignInRouter();

  const handlers = useHandlers();
  const classes = useStyles();

  return (
    <Container component="main" maxWidth="xs">
      <CssBaseline />
      <div className={classes.paper}>
        <Avatar className={classes.avatar}>
          <LockOutlinedIcon />
        </Avatar>
        <Typography component="h1" variant="h5">
          Sign in
        </Typography>
        <form className={classes.form} noValidate={true}>
          <TextField
            onChange={handlers.onChangeEmailForm}
            disabled={state.waitingSignIn}
            variant="outlined"
            margin="normal"
            required={true}
            fullWidth={true}
            id="email"
            label="Email Address"
            name="email"
            autoComplete="email"
            autoFocus={true}
          />
          <TextField
            onChange={handlers.onChangePasswordForm}
            disabled={state.waitingSignIn}
            variant="outlined"
            margin="normal"
            required={true}
            fullWidth={true}
            name="password"
            label="Password"
            type="password"
            id="password"
            autoComplete="current-password"
          />
          <FormControlLabel
            control={<Checkbox value="remember" color="primary" disabled={state.waitingSignIn}/>}
            label="Remember me"
          />
          <div className={classes.wrapper}>
            <Button
              id="submit-sign-in-request-button"
              disabled={state.waitingSignIn}
              onClick={handlers.onClickSignInSubmitButton}
              fullWidth={true}
              variant="contained"
              color="primary"
              className={classes.submit}
            >
              Sign In
            </Button>
            {state.waitingSignIn && <CircularProgress size={24} className={classes.buttonProgress} />}
          </div>
          <Grid container={true}>
            <Grid item={true} xs={true}>
              <Link href="#" variant="body2">
                Forgot password?
              </Link>
            </Grid>
            <Grid item={true}>
              <Link href="#" variant="body2">
                {"Don't have an account? Sign Up"}
              </Link>
            </Grid>
          </Grid>
        </form>
      </div>
      <Box mt={8}>
        <Copyright />
      </Box>
    </Container>
  );
};

export default SignIn;
