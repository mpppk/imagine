import {END, eventChannel} from "redux-saga";
import {wsActionCreators} from "../actions/ws";
import {AnyAction} from "typescript-fsa";
import {call, put, take} from "@redux-saga/core/effects";

let ws: WebSocket | null = null;

function wsChannel() {
  return eventChannel(emitter => {
      // FIXME
      ws = new WebSocket('ws://localhost:1323/ws')

      ws.onopen = function (event) {
        emitter(wsActionCreators.open(event))
      }

      ws.onclose = function (_event) {
        emitter(END)
      }

      ws.onerror = function (event) {
        emitter(wsActionCreators.error(event))
      }

      ws.onmessage = function (event) {
        emitter(wsActionCreators.message(event))
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

export function send(action: AnyAction): boolean {
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
      yield put(yield take(chan))
    }
  } finally {
    yield put(wsActionCreators.close())
  }
}