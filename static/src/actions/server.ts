import actionCreatorFactory from 'typescript-fsa';

const serverActionCreatorFactory = actionCreatorFactory('SERVER');

export const serverActionCreators = {
  startDirectoryScanning: serverActionCreatorFactory<void>(
    'START_DIRECTORY_SCANNING'
  ),
  scanningImages: serverActionCreatorFactory<string[]>(
    'SCANNING_IMAGES'
  ),
};
