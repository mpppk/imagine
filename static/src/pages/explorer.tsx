import * as React from 'react';
import {AssetTable} from "../components/AssetTable";
import {State} from "../reducers/reducer";
import {useDispatch, useSelector} from "react-redux";
import {assetActionCreators} from "../actions/asset";

// tslint:disable-next-line:variable-name
const Explorer: React.FC = () => {
  const globalState = useSelector((s: State) => s.global);
  const dispatch = useDispatch();

  const loadNextPage = () => {
    if (globalState.currentWorkSpace !== null) {
      dispatch(assetActionCreators.requestAssets({
        requestNum: 10,
        workSpaceName: globalState.currentWorkSpace.name,
      }));
    }
    return null;
  };

  return (
    <div>
      <AssetTable
        assets={globalState.assets}
        onRequestNextPage={loadNextPage}
        hasMoreAssets={globalState.hasMoreAssets}
        isScanningAssets={globalState.isScanningAssets}
        workspace={globalState.currentWorkSpace}
      />
    </div>
  );
}

export default Explorer;
