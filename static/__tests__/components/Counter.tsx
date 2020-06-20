// import Button from '@material-ui/core/Button/Button';
import { Button } from '@material-ui/core';
import { shallow } from 'enzyme';
import * as React from 'react';
import Counter from '../../src/components/Counter';

describe('Counter', () => {
  // tslint:disable-next-line no-empty
  const emptyButtonHandler = () => {};
  it('has 3 buttons', async () => {
    const wrapper = shallow(
      <Counter
        count={0}
        onClickIncrementButton={emptyButtonHandler}
        onClickDecrementButton={emptyButtonHandler}
        onClickIncrementLaterButton={emptyButtonHandler}
      />
    );
    expect(wrapper.find(Button)).toHaveLength(3);
  });
});
