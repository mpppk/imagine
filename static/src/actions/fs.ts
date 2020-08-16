import actionCreatorFactory from 'typescript-fsa';

const fsActionCreatorFactory = actionCreatorFactory('FS');

export interface WSPayload {
  workSpaceName: string
}

export const fsActionCreators = {
  scanCancel: fsActionCreatorFactory<WSPayload>('SCAN/CANCEL'),
  scanFinish: fsActionCreatorFactory<WSPayload>(
    'SCAN/FINISH'
  ),
  scanStart: fsActionCreatorFactory<WSPayload>(
    'SCAN/START'
  ),
};
