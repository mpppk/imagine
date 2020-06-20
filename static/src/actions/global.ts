import actionCreatorFactory from 'typescript-fsa';
import { User } from '../models/models';

const globalActionCreatorFactory = actionCreatorFactory('GLOBAL');

interface SignInRequest {
  email: string;
  password: string;
}

interface SignInResult {
  jwt: string; // FIXME
  user: User;
}

interface SignInError {
  error: Error;
}

export const globalActionCreators = {
  clickSignInSubmitButton: globalActionCreatorFactory<SignInRequest>(
    'CLICK_SIGN_IN_SUBMIT_BUTTON'
  ),
};

export const globalAsyncActionCreators = {
  signIn: globalActionCreatorFactory.async<
    SignInRequest,
    SignInResult,
    SignInError
  >('SIGN_IN'),
  signOut: globalActionCreatorFactory.async('SIGN_OUT'),
};
