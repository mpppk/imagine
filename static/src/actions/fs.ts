import actionCreatorFactory from 'typescript-fsa';

const fsActionCreatorFactory = actionCreatorFactory('FS');

export interface WSPayload {
  workSpaceName: string
}

interface FSScanRunningPayload extends WSPayload{
  paths: string[]
}

export const fsActionCreators = {
  scanCancel: fsActionCreatorFactory<WSPayload>('SCAN/CANCEL'),
  scanFinish: fsActionCreatorFactory<WSPayload>('SCAN/FINISH'),
  scanStart: fsActionCreatorFactory<WSPayload>('SCAN/START'),
  scanRunning: fsActionCreatorFactory<FSScanRunningPayload>('SCAN/RUNNING'),
};
