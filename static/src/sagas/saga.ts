import { all } from '@redux-saga/core/effects';
import { globalAsyncActionCreators } from '../actions/global';
import { watchIncrementAsync } from './counter';
import {
  signOutWorker,
  takeEveryAction,
  watchClickSignInSubmitButton,
} from './session';
import {connectToServer, wsSaga} from "./ws";
import {indexActionCreators} from "../actions";

export default function* rootSaga() {
  yield all([
    watchIncrementAsync(),
    watchClickSignInSubmitButton(),
    takeEveryAction(globalAsyncActionCreators.signOut.started, signOutWorker)(),
    wsSaga(),
    ...connectToServer(indexActionCreators),
  ]);
}
