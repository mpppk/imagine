import actionCreatorFactory from 'typescript-fsa';
import {WSPayload} from "./workspace";

const fsActionCreatorFactory = actionCreatorFactory('FS');

interface FSScanRunningPayload extends WSPayload{
  paths: string[]
}

export const fsActionCreators = {
  scanCancel: fsActionCreatorFactory<WSPayload>('SCAN/CANCEL'),
  scanFinish: fsActionCreatorFactory<WSPayload>('SCAN/FINISH'),
  scanStart: fsActionCreatorFactory<WSPayload>('SCAN/START'),
  scanRunning: fsActionCreatorFactory<FSScanRunningPayload>('SCAN/RUNNING'),
};
