import {END, eventChannel} from "redux-saga";
import {wsActionCreators} from "../actions/ws";
import {ActionCreator, AnyAction} from "typescript-fsa";
import {call, put, take} from "@redux-saga/core/effects";
import {takeEvery} from "redux-saga/effects";

let ws: WebSocket | null = null;

function wsChannel() {
  return eventChannel(emitter => {
      // FIXME
      ws = new WebSocket('ws://localhost:1323/ws')

      ws.onopen = function (event) {
        emitter({type: 'open', event})
      }

      ws.onclose = function (_event) {
        emitter(END)
      }

      ws.onerror = function (event) {
        emitter({type: 'error', event})
      }

      ws.onmessage = function (event) {
        emitter({type: 'message', event})
      }

      return () => {
        if (ws === null) {
          return
        }
        ws.close();
        ws = null;
      }
    }
  )
}

export function dispatchToServer(action: AnyAction): boolean {
  if (ws === null) {
    return false
  }
  ws.send(JSON.stringify(action));
  return true
}

export function* wsSaga() {
  const chan = yield call(wsChannel)
  try {
    while (true) {
      const e = yield take(chan);
      switch (e.type) {
        case 'open':
          yield put(wsActionCreators.open(e.event))
          break;
        case 'error':
          yield put(wsActionCreators.error(e.event))
          break;
        case 'message':
          yield put(wsActionCreators.message(e.event));
          yield put(JSON.parse(e.event.data));
      }
    }
  } finally {
    yield put(wsActionCreators.close())
  }
}

export function* pipeToServer(action: AnyAction) {
  yield takeEvery(action.type, (action) => dispatchToServer(action),
  );
}

type ActionCreators = {
  [key: string]: ActionCreator<any>
}

export function connectToServer(actions: ActionCreators) {
  const effects = [];
  for (let k of Object.keys(actions)) {
    effects.push(pipeToServer(actions[k]))
  }
  return effects
}
