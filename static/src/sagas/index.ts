import {takeEvery} from 'redux-saga/effects';
import {indexActionCreators} from "../actions";
import {send} from "./ws";

export function* watchClickAddDirectoryButton() {
  yield takeEvery(
    indexActionCreators.clickAddDirectoryButton.type,
    (action) => send(action),
  );
}
