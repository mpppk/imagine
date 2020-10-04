import {Theme} from "@material-ui/core";
import {makeStyles} from "@material-ui/core/styles";
import React, {useMemo} from "react";
import {AssetWithIndex} from "../models/models";
import {isDefaultBox} from "../util";

const useStyles = makeStyles((theme: Theme) => {
  return {
    rect: {
      fill: 'transparent',
      stroke: theme.palette.primary.light,
    }
  }
});

interface Props {
  src: string
  asset: AssetWithIndex
}

// const useHandlers = (props: Props) => {
//   return useMemo(() => {
//     return {}
//   }, [props]);
// }
//

type Classes = ReturnType<typeof useStyles>;
const useViewState = (props: Props, classes: Classes) => {
  return useMemo(() => {
    const boxes = props.asset.boundingBoxes ?? [];
    const rectProps = boxes.map((box) => {
      const rectProp = {
        className: classes.rect,
        width: box.width,
        height: box.height,
        x: box.x,
        y: box.y,
        key: box.id,
      };
      return isDefaultBox(box) ? {...rectProp, width: '100%', height: '100%'} : rectProp;
    });
    return {rectProps};
  }, [props, classes]);
}

// tslint:disable-next-line:variable-name
export const ImagePreview: React.FC<Props> = (props) => {
  const classes = useStyles();
  const viewState = useViewState(props, classes);
  // const handlers = useHandlers(props);

  return (
    <div>
      <svg id="canvas" viewBox="0 0 500 500" width="500" height="500">
        <image href={props.src} width={'100%'}/>
        {viewState.rectProps.map((rectProp) => {
          return <rect {...rectProp} key={rectProp.key}/>
        })}
        {/*<rect className={classes.rect} width={props.asset.} height="100" x="0" y="0"/>*/}
      </svg>
    </div>
  )
}
