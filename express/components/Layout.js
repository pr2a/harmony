import Head from 'next/head';
import './queries.css';
import './style.css';

export default ({ children, title = 'Lottery' }) => (
  <div>
    <Head>
      <meta charSet="UTF-8" />
      <meta name="viewport" content="width=device-width, initial-scale=1.0" />

      <link
        href="https://fonts.googleapis.com/css?family=Quicksand:400"
        rel="stylesheet"
      />
      <link
        href="https://fonts.googleapis.com/css?family=Open+Sans:400"
        rel="stylesheet"
      />

      <title>Lottery | ...</title>
    </Head>
    <header className="header">
      <div className="header__logoBox">
        <img
          src="static/img/logo.svg"
          alt="Lottery logo"
          className="header__logo"
        />
      </div>
    </header>
    {children}
  </div>
);
