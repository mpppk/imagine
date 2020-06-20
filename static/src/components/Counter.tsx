import Button from '@material-ui/core/Button/Button';
import Typography from '@material-ui/core/Typography/Typography';
import * as React from 'react';

export interface ICounterProps {
  count: number;
  onClickIncrementButton: (_: void) => void;
  onClickDecrementButton: (_: void) => void;
  onClickIncrementLaterButton: (_: void) => void;
}

type ClickEvent = React.MouseEvent<HTMLElement, MouseEvent>;

export default (props: ICounterProps) => {
  const handleClickIncrementButton = (_e: ClickEvent) => {
    props.onClickIncrementButton();
  };
  const handleClickDecrementButton = (_e: ClickEvent) => {
    props.onClickDecrementButton();
  };
  const handleClickIncrementLaterButton = (_e: ClickEvent) => {
    props.onClickIncrementLaterButton();
  };

  return (
    <div>
      <Typography variant="h4" gutterBottom={true}>
        Count: <span>{props.count}</span>
      </Typography>
      <Button
        onClick={handleClickIncrementButton}
        variant="contained"
        color="primary"
      >
        +1
      </Button>
      <Button
        onClick={handleClickDecrementButton}
        variant="contained"
        color="primary"
      >
        -1
      </Button>
      <Button
        onClick={handleClickIncrementLaterButton}
        variant="contained"
        color="primary"
      >
        +1 later
      </Button>
    </div>
  );
};
