import Button from '@material-ui/core/Button';
import Link from 'next/link';
import React from 'react';

// tslint:disable-next-line variable-name
const LoginButton: React.FunctionComponent = () => {
  return (
    <Link href={'/signin'}>
      <Button color="inherit">Sign In</Button>
    </Link>
  );
};

export default LoginButton;
