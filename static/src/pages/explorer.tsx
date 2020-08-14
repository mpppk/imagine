import * as React from 'react';
import {AssetTable} from '../components/AssetTable';
import {useVirtualizedAsset} from "../hooks";

// tslint:disable-next-line:variable-name
const Explorer: React.FC = () => {
  const virtualizedAssetProps = useVirtualizedAsset();

  return (
    <div>
      <AssetTable
        {...virtualizedAssetProps}
      />
    </div>
  );
}

export default Explorer;
