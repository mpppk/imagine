import { useMemo } from 'react'
import {useDispatch, useSelector} from 'react-redux'
import { ActionCreatorsMapObject, bindActionCreators } from 'redux';
import {State} from "./reducers/reducer";
import {assetActionCreators} from "./actions/asset";

export function useActions<T extends ActionCreatorsMapObject>(actions: T, deps?: any[]) {
  const dispatch = useDispatch()
  return useMemo(
    () => {
      return bindActionCreators(actions, dispatch)
    },
    deps ? [dispatch, ...deps] : [dispatch]
  )
}

export function useVirtualizedAsset() {
  const globalState = useSelector((s: State) => s.global);
  const dispatch = useDispatch();

  const loadNextPage = () => {
    if (globalState.currentWorkSpace !== null) {
      dispatch(assetActionCreators.scanRequest({
        requestNum: 10,
        workSpaceName: globalState.currentWorkSpace.name,
      }));
    }
    return null;
  };
  return {
    assets: globalState.assets,
    onRequestNextPage: loadNextPage,
    hasMoreAssets: globalState.hasMoreAssets,
    isScanningAssets: globalState.isScanningAssets,
    workspace: globalState.currentWorkSpace,
  };
}