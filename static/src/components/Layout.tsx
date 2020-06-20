import Head from 'next/head';
import Link from 'next/link';
import * as React from 'react';

interface IProps {
  title?: string;
}

// tslint:disable-next-line variable-name
const Layout: React.FunctionComponent<IProps> = ({
  children,
  title = 'This is the default title'
}) => (
  <div>
    <Head>
      <title>{title}</title>
      <meta charSet="utf-8" />
      <meta name="viewport" content="initial-scale=1.0, width=device-width" />
    </Head>
    <header>
      <nav>
        <Link href="/">
          <a>Home</a>
        </Link>{' '}
        |{' '}
        <Link href="/about">
          <a>About</a>
        </Link>
      </nav>
    </header>
    {children}
    <footer>I'm here to stay</footer>
  </div>
);

export default Layout;
