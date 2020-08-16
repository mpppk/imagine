import actionCreatorFactory from 'typescript-fsa';
import {Asset} from "../models/models";
import {WSPayload} from "./workspace";

const assetActionCreatorFactory = actionCreatorFactory('ASSET');

interface RequestAssetsPayload extends WSPayload {
  requestNum: number
}

interface ScanRunningPayload extends WSPayload{
  assets: Asset[]
}

export const assetActionCreators = {
  requestAssets: assetActionCreatorFactory<RequestAssetsPayload>('REQUEST_ASSETS'),
  scanRunning: assetActionCreatorFactory<ScanRunningPayload>('SCAN/RUNNING'),
  scanFinish: assetActionCreatorFactory<WSPayload>('SCAN/FINISH'),
};

