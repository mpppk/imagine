import actionCreatorFactory from 'typescript-fsa';

const serverActionCreatorFactory = actionCreatorFactory('SERVER');

export interface WSPayload {
  workSpaceName: string
}

interface ScanningImagesPayload extends WSPayload{
  paths: string[]
}

export const serverActionCreators = {
  startDirectoryScanning: serverActionCreatorFactory<WSPayload>(
    'START_DIRECTORY_SCANNING'
  ),
  cancelDirectoryScanning: serverActionCreatorFactory<WSPayload>(
    'CANCEL_DIRECTORY_SCANNING'
  ),
  finishDirectoryScanning: serverActionCreatorFactory<WSPayload>(
    'FINISH_DIRECTORY_SCANNING'
  ),
  scanningImages: serverActionCreatorFactory<ScanningImagesPayload>(
    'SCANNING_IMAGES'
  ),
};
