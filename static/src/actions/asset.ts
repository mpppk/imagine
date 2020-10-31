import actionCreatorFactory from 'typescript-fsa';
import {Asset, Query} from "../models/models";
import {WSPayload} from "./workspace";

const assetActionCreatorFactory = actionCreatorFactory('ASSET');

interface AssetScanRequest extends WSPayload {
  requestNum: number
  queries: Query[]
  reset: boolean
}

interface ScanRunningPayload extends WSPayload{
  assets: Asset[]
}

export const assetActionCreators = {
  scanRequest: assetActionCreatorFactory<AssetScanRequest>('SCAN/REQUEST'),
  scanRunning: assetActionCreatorFactory<ScanRunningPayload>('SCAN/RUNNING'),
  scanFinish: assetActionCreatorFactory<WSPayload>('SCAN/FINISH'),
};

