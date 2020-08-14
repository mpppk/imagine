import actionCreatorFactory from 'typescript-fsa';
import {WSPayload} from "./server";

const assetActionCreatorFactory = actionCreatorFactory('ASSET');

interface RequestAssetsPayload extends WSPayload {
  requestNum: number
}

export const assetActionCreators = {
  requestAssets: assetActionCreatorFactory<RequestAssetsPayload>('REQUEST_ASSETS'),
};

