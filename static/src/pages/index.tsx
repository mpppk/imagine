import Typography from '@material-ui/core/Typography/Typography';
import { NextPage } from 'next';
import React from 'react';
import { useSelector } from 'react-redux';
import { counterActionCreators } from '../actions/counter';
import Counter from '../components/Counter';
import { useActions } from '../hooks';
import { State } from '../reducers/reducer';

// tslint:disable-next-line variable-name
export const Index: NextPage = () => {
  const handlers = useActions(counterActionCreators);
  const count = useSelector((state: State) => state.counter.count);

  return (
    <div>
      <Typography variant="h2" gutterBottom={true}>
        Counter sample
      </Typography>
      <Counter
        count={count}
        onClickIncrementButton={handlers.clickIncrementButton}
        onClickDecrementButton={handlers.clickDecrementButton}
        onClickIncrementLaterButton={handlers.clickAsyncIncrementButton}
      />
    </div>
  );
};

export default Index;
