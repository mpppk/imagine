import actionCreatorFactory from 'typescript-fsa';
import { WSPayload } from './workspace';

const fsActionCreatorFactory = actionCreatorFactory('FS');

export interface FSScanRequestPayload extends WSPayload {
  basePath: string;
}

interface FSScanRunningPayload extends WSPayload {
  foundedAssetsNum: number;
}

export interface FSScanStartPayload extends WSPayload {
  basePath: string;
}

export interface BaseDirSelectPayload extends WSPayload {
  basePath: string;
}

export const fsActionCreators = {
  scanRequest: fsActionCreatorFactory<FSScanRequestPayload>('SCAN/REQUEST'),
  scanCancel: fsActionCreatorFactory<WSPayload>('SCAN/CANCEL'),
  scanFinish: fsActionCreatorFactory<WSPayload>('SCAN/FINISH'),
  scanStart: fsActionCreatorFactory<FSScanStartPayload>('SCAN/START'),
  scanRunning: fsActionCreatorFactory<FSScanRunningPayload>('SCAN/RUNNING'),
  baseDirSelect: fsActionCreatorFactory<BaseDirSelectPayload>(
    'BASE_DIR/SELECT'
  ),
};
